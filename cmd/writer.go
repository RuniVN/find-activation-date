package cmd

import (
	"encoding/csv"
	"fmt"
	"io"

	"golang.org/x/sync/syncmap"
)

// createFile creates a file
//
// createFile receives a fileSystem that used to mock later
func createFile(fs fileSystem, fileName string) (file, error) {
	return fs.Create(fileName)
}

// writeToCSV takes a map, parse this map then write to the file
func writeToCSV(fi file, ms *syncmap.Map) error {
	writer := csv.NewWriter(fi)
	defer writer.Flush()

	writer.Write([]string{PhoneNumberText, RealActivationDateText})

	ms.Range(func(k, v interface{}) bool {
		row, ok := v.(*Result)
		if !ok {
			return false
		}
		err := writer.Write([]string{k.(string), row.RealActivationDate})
		if err != nil {
			return false
		}

		return true
	})

	return nil
}

// WriteRow takes a row then write into a writer
func WriteRow(f io.Writer, row Row) error {
	_, err := f.Write([]byte(fmt.Sprintf("%s,%s\n", row.ActivationDate, row.DeactivationDate)))
	if err != nil {
		return err
	}

	return nil
}
