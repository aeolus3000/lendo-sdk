package banking

import "testing"

func Test_checkStatus(t *testing.T) {
	type args struct {
		status string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"pending",
			args{status: "pending"},
			false,
		}, {
			"completed",
			args{status: "completed"},
			false,
		}, {
			"rejected",
			args{status: "rejected"},
			false,
		}, {
			"something else",
			args{status: "something else"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkStatus(tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("checkStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
