package cmd

import (
	"encoding/csv"
	"os"
	"testing"

	"golang.org/x/sync/syncmap"
)

func TestLessRow(t *testing.T) {
	type args struct {
		row1 Row
		row2 Row
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Less",
			args: args{
				row1: Row{ActivationDate: "2015-02-02", DeactivationDate: "2015-05-02"},
				row2: Row{ActivationDate: "2015-08-02", DeactivationDate: "2015-12-02"},
			},
			want: true,
		},
		{
			name: "More",
			args: args{
				row1: Row{ActivationDate: "2015-02-02", DeactivationDate: "2015-05-02"},
				row2: Row{ActivationDate: "2014-08-02", DeactivationDate: "2014-09-02"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lessRow(tt.args.row1, tt.args.row2); got != tt.want {
				t.Errorf("lessRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExportCSV(t *testing.T) {
	key := "098812345"
	result := &Result{"2015-02-02"}

	ms := &syncmap.Map{}
	ms.Store(key, result)

	err := ExportCSV(ms)
	if err != nil {
		t.Fatalf("ExportCSV() failed, err = %v", err)
	}

	f, err := os.Open(ResultFile)
	if err != nil {
		t.Fatalf("failed to open file %s, err = %v", ResultFile, err)
	}

	r := csv.NewReader(f)

	r1, err := r.Read()
	if err != nil {
		t.Fatalf("read temp file failed, err = %v", err)
	}

	r2, err := r.Read()
	if err != nil {
		t.Fatalf("read temp file failed, err = %v", err)
	}

	if r1[0] != PhoneNumberText || r1[1] != RealActivationDateText {
		t.Errorf("data written not expected, expected = %s %s, got %s %s", PhoneNumberText, RealActivationDateText, r1[0], r1[1])
	}

	if r2[0] != key || r2[1] != result.RealActivationDate {
		t.Errorf("data written not expected, expected = %s %s, got %s %s", key, result.RealActivationDate, r2[0], r2[1])
	}

	err = os.Remove(ResultFile)
	if err != nil {
		t.Fatalf("failed to clean up, err = %v", err)
	}
}

func TestGetAllRowsOneFile(t *testing.T) {
	tempFile := "tmp.csv"
	f, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("failed to create file %s, err = %v", tempFile, err)
	}

	f.Write([]byte("2014-02-02, 2015-02-02"))

	rows, err := getAllRowsOneFile(tempFile)
	if err != nil {
		t.Fatalf("getAllRowsOneFile() failed, err = %v", err)
	}

	expected := 1
	if len(rows) != expected {
		t.Errorf("getAllRowsOneFile() failed, expected length = %d, got %d", expected, len(rows))
	}
}

func TestFindTheLatestActivation(t *testing.T) {
	realActivationDate := "2015-05-02"
	rows := []Row{
		{
			ActivationDate:   "2015-02-02",
			DeactivationDate: "2015-03-02",
		},
		{
			ActivationDate:   "2015-03-02",
			DeactivationDate: "2015-04-02",
		},
		{
			ActivationDate:   "2015-05-02",
			DeactivationDate: "",
		},
	}

	res := findTheLatestActivation(rows)

	if res.RealActivationDate != realActivationDate {
		t.Errorf("wrong result, expected %s, got %s", realActivationDate, res.RealActivationDate)
	}
}

func TestTearDown(t *testing.T) {
	err := TearDown()
	if err != nil {
		t.Errorf("failed to tear down, err = %v", err)
	}
}
