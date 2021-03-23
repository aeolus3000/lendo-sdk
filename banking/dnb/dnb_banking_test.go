// +build integration_test

package dnb

//These tests need the following container running:
//docker run -p 8000:8000 lendoab/interview-service:stable

import (
	"github.com/aeolus3000/lendo-sdk/banking"
	"github.com/google/uuid"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	maxProcessingTime = 21 * time.Second
)

var (
	dnbApi banking.BankingApi
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	dnbApi = NewDnbBanking(NewDnbDefaultConfiguration())
}

func shutdown() {

}

func TestCreateApplicationValid(t *testing.T) {
	expectedStatus := "pending"
	application := banking.Application{
		Id:        uuid.NewString(),
		FirstName: "Simon",
		LastName:  "Kopp",
	}
	applicationStatus, err := dnbApi.Create(&application)
	if err != nil {
		t.Errorf("Unexpected error creating application: %v", err)
	}
	if applicationStatus.Status != expectedStatus {
		t.Errorf("Application does not have \"%v\" status", expectedStatus)
	}
	if applicationStatus.Id != application.Id {
		t.Errorf("Response application id didn't match; got: %v; want: %v", applicationStatus.Id, application.Id)
	}
}

func TestCreateApplicationInValidUuid(t *testing.T) {
	expectedError := "duplicate application Id: {\"error\":\"the request payload is not valid.\"}"
	application := banking.Application{
		Id:        "invalid-uuid",
		FirstName: "Simon",
		LastName:  "Kopp",
	}
	_, err := dnbApi.Create(&application)
	if err == nil {
		t.Errorf("Expected error when creating an application with an invalid uuid")
	}
	if err != nil && err.Error() != expectedError {
		t.Errorf("Expected error message didn't equal; got: %v; want: %v", err.Error(), expectedError)
	}
}

func TestCreateApplicationDuplicateUuid(t *testing.T) {
	expectedError := "duplicate application Id: {\"error\":\"the application ID is already used\"}"
	application := banking.Application{
		Id:        uuid.NewString(),
		FirstName: "Simon",
		LastName:  "Kopp",
	}
	_, _ = dnbApi.Create(&application)
	// create a second time to see duplicate error
	_, err := dnbApi.Create(&application)
	if err == nil {
		t.Errorf("Expected error when creating an application with a duplicate uuid")
	}
	if err != nil && err.Error() != expectedError {
		t.Errorf("Expected error message didn't equal; got: %v; want: %v", err.Error(), expectedError)
	}
}

func TestCheckStatusValidUuid(t *testing.T) {
	expectedStatus := "completedrejected"
	application := banking.Application{
		Id:        uuid.NewString(),
		FirstName: "Simon",
		LastName:  "Kopp",
	}
	_, _ = dnbApi.Create(&application)
	time.Sleep(maxProcessingTime)
	applicationStatus, err := dnbApi.CheckStatus(application.Id)
	if err != nil {
		t.Errorf("Checking the status failed unexpectedly: %v", err)
	}
	if !strings.Contains(expectedStatus, applicationStatus.Status) {
		t.Errorf("The application status didn't match the expactation; got = %v; want one of = %v",
			applicationStatus.Status, expectedStatus)
	}
	if applicationStatus.Id != application.Id {
		t.Errorf("Response application id didn't match; got: %v; want: %v", applicationStatus.Id, application.Id)
	}
	if applicationStatus.JobId == "" {
		t.Errorf("Didn't receive a Job ID from banking service")
	}
}

func TestCheckStatusInvalidUuid(t *testing.T) {
	tests := []struct {
		name          string
		uuid          string
		expectedError string
	}{
		{
			"Missing uuid",
			"",
			"invalid uuid: {\"error\":\"application_id is missing\"}",
		}, {
			"invalid uuid",
			"this-is-not-a-valid-uuid",
			"invalid uuid: {\"error\":\"application_id is not a valid uuid\"}",
		}, {
			"not existing uuid",
			uuid.NewString(),
			"uuid does not exist: {\"error\":\"failed to get job, the application_id doesn't exist.\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dnbApi.CheckStatus(tt.uuid)
			if err == nil {
				t.Errorf("Expected error when requesting check status with an invalid uuid")
			}
			if err != nil && err.Error() != tt.expectedError {
				t.Errorf("Expected error message didn't equal; got: %v; want: %v", err.Error(), tt.expectedError)
			}
		})
	}
}
