package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/Sirupsen/logrus"

	"golang.org/x/sync/syncmap"
)

// define vars
var (
	Layout               = "2006-01-02"        // Layout for datetime format csv
	ResultFile           = "result.csv"        // result file name
	TempFolder           = "./tmp"             // folder temp
	PhoneNumberText      = "PHONE_NUMBER"      // text of phone number stored in result file
	ReactivationDateText = "REACTIVATION_DATE" // text of reactivation date stored in result file
)

// osFS implements fileSystem using the local disk.
type osFS struct{}

// Create is overrided function from file system interface
func (osFS) Create(name string) (file, error) { return os.Create(name) }

// ExportCSV receives a sync map then write result into csv file
func ExportCSV(ms *syncmap.Map) error {
	lf := logrus.Fields{"func": "cmd.ExportCSV"}

	var fs fileSystem = osFS{}

	fi, err := createFile(fs, ResultFile)
	if err != nil {
		logrus.WithFields(lf).WithError(err).Error("failed to create file")
		return err
	}
	defer fi.Close()

	return writeToCSV(fi, ms)
}

// Preprocessing separtes all data of csv file into multiple files
//
// each file has the name which is the phone number
// each file contains activation date and deactivation date of one number
func Preprocessing(f io.Reader) ([]string, error) {
	lf := logrus.Fields{"func": "cmd.Preprocessing"}

	var lFiles []string
	r := csv.NewReader(f)

	var fs fileSystem = osFS{}
	var mf = make(map[string]file)

	r.Read()
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.WithFields(lf).WithError(err).Error("failed to read csv file")
			return nil, err
		}

		var fi file
		var ok bool
		if fi, ok = mf[record[0]]; !ok {
			lFiles = append(lFiles, record[0])
			fi, err = createFile(fs, fmt.Sprintf("%s/%s", TempFolder, record[0]))
			if err != nil {
				logrus.WithFields(lf).WithError(err).Error("failed to read csv file")
				return nil, err
			}
			mf[record[0]] = fi
		}

		err = WriteRow(fi, Row{ActivationDate: record[1], DeactivationDate: record[2]})
		if err != nil {
			logrus.WithFields(lf).WithError(err).Error("failed to write row")
			return nil, err
		}
	}

	for _, v := range mf {
		v.Close()
	}

	return lFiles, nil
}

// ProcessOneFile processes only one file csv
//
// Return value includes the result and error which result is
// real activation date
func ProcessOneFile(fileName string) (*Result, error) {
	lf := logrus.Fields{"func": "cmd.ProcessOneFile"}

	rows, err := getAllRowsOneFile(fileName)
	if err != nil {
		logrus.WithFields(lf).WithError(err).Error("failed to get all rows of one file")
		return nil, err
	}
	return findTheLatestActivation(rows), nil
}

// getAllRowsOneFile returns all rows in single file, merged
// into a []Row model
func getAllRowsOneFile(fileName string) ([]Row, error) {
	lf := logrus.Fields{"func": "cmd.ProcessOneFile"}

	f, err := os.Open(fmt.Sprintf("%s/%s", TempFolder, fileName))
	if err != nil {
		logrus.WithFields(lf).WithError(err).Errorf("failed to open file %s", fileName)
		return nil, err
	}
	defer f.Close()

	var rows []Row
	r := csv.NewReader(f)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.WithFields(lf).WithError(err).Error("failed to read csv file")
			return nil, err
		}

		rows = append(rows, Row{ActivationDate: record[0], DeactivationDate: record[1]})
	}

	return rows, nil
}

// findTheLatestActivation does 2 things
//
// 1. Sort the file by datetime asc
// 2. Traceback from the latest row, and find the row
// which activation date != deactivation date of the previous row
func findTheLatestActivation(rows []Row) *Result {
	sort.Slice(rows, func(i, j int) bool {
		return lessRow(rows[i], rows[j])
	})

	if len(rows) < 1 {
		return &Result{}
	}

	if len(rows) == 1 {
		return &Result{rows[0].ActivationDate}
	}

	i := len(rows) - 1
	for i > 1 {
		if rows[i].ActivationDate != rows[i-1].DeactivationDate {
			return &Result{rows[i].ActivationDate}
		}
		i--
	}

	return &Result{rows[i-1].ActivationDate}
}

// lessRow returns whether row1 is before row2 in timeline
func lessRow(row1, row2 Row) bool {
	lf := logrus.Fields{"func": "cmd.lessRow"}

	t1, err := time.Parse(Layout, row1.ActivationDate)
	if err != nil && row1.ActivationDate != "" {
		logrus.WithFields(lf).WithError(err).Errorf("failed to parse date time %s, skip", row1.ActivationDate)
		return true
	}

	t2, err := time.Parse(Layout, row2.DeactivationDate)
	if err != nil && row2.DeactivationDate != "" {
		logrus.WithFields(lf).WithError(err).Errorf("failed to parse date time %s, skip", row2.DeactivationDate)
		return true
	}

	return t1.Before(t2)
}

// TearDown cleans up all temporary files
func TearDown() error {
	lf := logrus.Fields{"func": "cmd.TearDown"}

	err := os.RemoveAll(TempFolder)
	if err != nil {
		logrus.WithFields(lf).WithError(err).Errorf("failed to remove folder %s", TempFolder)
		return err
	}
	err = os.MkdirAll(TempFolder, os.ModeDir)
	if err != nil {
		logrus.WithFields(lf).WithError(err).Errorf("failed to mkdir folder %s", TempFolder)
		return err
	}

	return nil
}
