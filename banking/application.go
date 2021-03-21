package banking

import "fmt"

const (
	STATUS_PENDING = "pending"
	STATUS_COMPLETED = "completed"
	STATUS_REJECTED = "rejected"
)

type Application struct {
	Id        string
	Firstname string
	Lastname  string
	Status    string
	JobId	  string
}

func checkStatus(status string) error {
	switch status {
	case STATUS_PENDING, STATUS_COMPLETED, STATUS_REJECTED:
		return nil
	default:
		return fmt.Errorf("invalid status: %v", status)
	}
}