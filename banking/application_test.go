package banking

import (
	"bytes"
	"google.golang.org/protobuf/proto"
	"reflect"
	"testing"
)

func Test_applicationProcessed(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{
			"rejected",
			"rejected",
			true,
		}, {
			"completed",
			"completed",
			true,
		}, {
			"pending",
			"pending",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApplicationProcessed(tt.status); got != tt.want {
				t.Errorf("applicationProcessed() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func Test_DeserializeToApplication(t *testing.T) {
	expectedApplication := Application{
		Id:        "8cb9fd11-fc05-45d2-93fa-23f62dee24b2",
		FirstName: "Employer",
		LastName:  "Of the month",
		Status:    "pending",
		JobId:     "8cb9fd11-fc05-45d2-93fa-23f62dee24b2",
	}
	applicationBytes := bytes.NewBuffer([]byte{10, 36, 56, 99, 98, 57, 102, 100, 49, 49, 45, 102, 99, 48, 53, 45, 52, 53, 100, 50, 45, 57,
		51, 102, 97, 45, 50, 51, 102, 54, 50, 100, 101, 101, 50, 52, 98, 50, 18, 8, 69, 109, 112, 108, 111, 121, 101, 114, 26, 12, 79, 102, 32, 116,
		104, 101, 32, 109, 111, 110, 116, 104, 34, 7, 112, 101, 110, 100, 105, 110, 103, 42, 36, 56, 99, 98, 57, 102, 100, 49, 49, 45, 102, 99, 48,
		53, 45, 52, 53, 100, 50, 45, 57, 51, 102, 97, 45, 50, 51, 102, 54, 50, 100, 101, 101, 50, 52, 98, 50})
	applicationBytes2 := bytes.NewBuffer([]byte{10, 36, 56, 99, 98, 57, 102, 100, 49, 49, 45, 102, 99, 48, 53, 45, 52, 53, 100, 50, 45, 57, 51, 102, 97, 45, 50, 51, 102,
		54, 50, 100, 101, 101, 50, 52, 98, 50, 18, 8, 69, 109, 112, 108, 111, 121, 101, 114, 26, 12, 79, 102, 32, 116, 104, 101, 32, 109, 111, 110,
		116, 104, 34, 7, 80, 69, 78, 68, 73, 78, 71, 42, 36, 56, 99, 98, 57, 102, 100, 49, 49, 45, 102, 99, 48, 53, 45, 52, 53, 100, 50, 45, 57, 51, 102,
		97, 45, 50, 51, 102, 54, 50, 100, 101, 101, 50, 52, 98, 50}) //status "PENDING" in capital letters -> deserialize should lower case
	tests := []struct {
		name       string
		byteBuffer *bytes.Buffer
		wantErr    bool
	}{
		{
			"no error",
			applicationBytes,
			false,
		}, {
			"no error, should make to lower on status",
			applicationBytes2,
			false,
		}, {
			"error",
			bytes.NewBuffer([]byte{10, 36, 56, 99, 98, 57, 102}),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			application, err := DeserializeToApplication(tt.byteBuffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("deserializeToApplication() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && application != nil {
				t.Error("deserializeToApplication() application wasn't empty in error case")
			}
			if err == nil && !proto.Equal(application, &expectedApplication) {
				t.Error("deserializeToApplication() applications didn't match")
			}
		})
	}
}

func Test_SerializeToApplication(t *testing.T) {
	application := Application{
		Id:        "8cb9fd11-fc05-45d2-93fa-23f62dee24b2",
		FirstName: "Employer",
		LastName:  "Of the month",
		Status:    "pending",
		JobId:     "8cb9fd11-fc05-45d2-93fa-23f62dee24b2",
	}
	application2 := Application{
		Id:        "8cb9fd11-fc05-45d2-93fa-23f62dee24b2",
		FirstName: "Employer",
		LastName:  "Of the month",
		Status:    "PENDING",
		JobId:     "8cb9fd11-fc05-45d2-93fa-23f62dee24b2",
	}
	expectedApplicationBytes := bytes.NewBuffer([]byte{10, 36, 56, 99, 98, 57, 102, 100, 49, 49, 45, 102, 99, 48, 53, 45, 52, 53, 100, 50,
		45, 57, 51, 102, 97, 45, 50, 51, 102, 54, 50, 100, 101, 101, 50, 52, 98, 50, 18, 8, 69, 109, 112, 108, 111, 121, 101, 114, 26, 12, 79, 102, 32,
		116, 104, 101, 32, 109, 111, 110, 116, 104, 34, 7, 112, 101, 110, 100, 105, 110, 103, 42, 36, 56, 99, 98, 57, 102, 100, 49, 49, 45, 102, 99,
		48, 53, 45, 52, 53, 100, 50, 45, 57, 51, 102, 97, 45, 50, 51, 102, 54, 50, 100, 101, 101, 50, 52, 98, 50})
	tests := []struct {
		name        string
		application *Application
		wantErr     bool
	}{
		{
			"no error",
			&application,
			false,
		}, {
			"no error, should make to lower on status",
			&application2,
			false,
		},
		// Serializing error can't be simulated
		//{
		//	"error",
		//	???,
		//	true,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesBuffer, err := SerializeFromApplication(tt.application)
			if (err != nil) != tt.wantErr {
				t.Errorf("SerializeToApplication() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && bytesBuffer != nil {
				t.Error("SerializeToApplication() produced error but serializing result wasn't nil")
			}
			if err == nil && !reflect.DeepEqual(bytesBuffer, expectedApplicationBytes) {
				t.Error("SerializeToApplication() applications didn't match")
			}
		})
	}
}
