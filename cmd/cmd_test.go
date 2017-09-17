package cmd

import (
	"testing"
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
