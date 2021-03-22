package banking

import (
	"bytes"
	"fmt"
	"google.golang.org/protobuf/proto"
	"strings"
)

const (
	STATUS_PENDING   = "pending"
	STATUS_COMPLETED = "completed"
	STATUS_REJECTED  = "rejected"
)

func checkStatus(status string) error {
	switch status {
	case STATUS_PENDING, STATUS_COMPLETED, STATUS_REJECTED:
		return nil
	default:
		return fmt.Errorf("invalid status: %v", status)
	}
}

func ApplicationProcessed(status string) bool {
	if strings.Contains(status, STATUS_COMPLETED) ||
		strings.Contains(status, STATUS_REJECTED) {
		return true
	}
	return false
}

func DeserializeToApplication(protoBytesBuffer *bytes.Buffer) (*Application, error) {
	application := Application{}
	err := proto.Unmarshal(protoBytesBuffer.Bytes(), &application)
	if err != nil {
		return nil, fmt.Errorf("can not unmarshal application from bytes: %v", err)
	}
	application.Status = strings.ToLower(application.Status)
	return &application, nil
}

func SerializeFromApplication(application *Application) (*bytes.Buffer, error) {
	application.Status = strings.ToLower(application.Status)
	applicationBytes, err := proto.Marshal(application)
	if err != nil {
		return nil, fmt.Errorf("can not serialize appication to bytes: %v", err)
	}
	return bytes.NewBuffer(applicationBytes), nil
}