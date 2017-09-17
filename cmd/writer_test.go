package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

var (
	ErrWrite        = errors.New("failed to write")
	ErrCreate       = errors.New("failed to create")
	ErrNoPermission = errors.New("no permission")
)

type MockWriter struct {
	SimulateWriteError bool
	Bs                 []byte
}

func (mr *MockWriter) Write(bs []byte) (int, error) {
	if mr.SimulateWriteError {
		return 0, ErrWrite
	}

	mr.Bs = bs
	return len(bs), nil
}

func TestWriteRow(t *testing.T) {
	type args struct {
		f   io.Writer
		row Row
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful",
			args: args{
				f:   &MockWriter{},
				row: Row{ActivationDate: "2014-02-02", DeactivationDate: "2015-01-01"},
			},
			wantErr: false,
		}, {
			name: "Fail",
			args: args{
				f: &MockWriter{
					SimulateWriteError: true,
				},
				row: Row{ActivationDate: "2014-02-02", DeactivationDate: "2015-01-01"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WriteRow(tt.args.f, tt.args.row)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				}
				if err != ErrWrite {
					t.Errorf("Write() error = %v, expected error =  %v", err, ErrWrite)
				}
			}

			m := tt.args.f.(*MockWriter)
			var expected []byte
			if err == nil {
				expected = []byte(fmt.Sprintf("%s,%s\n", tt.args.row.ActivationDate, tt.args.row.DeactivationDate))
			}
			actual := m.Bs

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("data not expected, expected %+v got %+v", []byte(expected), actual)
			}

		})
	}
}

type MockFileSystem struct {
	osFS
	SimulateCreateError bool
}

type MockFile struct {
	file
}

type MockFileInfo struct {
	os.FileInfo
}

func (mf *MockFileSystem) Create(fileName string) (file, error) {
	if mf.SimulateCreateError {
		return nil, ErrCreate
	}

	return MockFile{}, nil
}

func TestCreateFile(t *testing.T) {
	type args struct {
		fs       fileSystem
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    file
		wantErr bool
	}{
		{
			name: "Successful",
			args: args{
				fs: &MockFileSystem{},
			},
			want:    MockFile{},
			wantErr: false,
		},
		{
			name: "CheckCreateError",
			args: args{
				fs: &MockFileSystem{
					SimulateCreateError: true,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createFile(tt.args.fs, tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("createFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
