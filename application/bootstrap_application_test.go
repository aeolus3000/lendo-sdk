package application

import (
	"lendo-sdk/configuration"
	"os"
	"testing"
)

type TestApplication struct {
	AbstractApplication
	sequence string
}

func (ta *TestApplication) Initialize(_ configuration.Configuration) {
	ta.sequence += "init"
}

func (ta *TestApplication) Execute() {
	ta.sequence += "execute"
}

func (ta *TestApplication) Shutdown() {
	ta.sequence += "update"
}

func TestNewBootstrapApplication(t *testing.T) {
	expectedSequence := "initexecuteupdate"
	testApplication := TestApplication{}
	NewBootstrapApplication(&testApplication, "testapplication").Execute(os.Args)
	if testApplication.sequence != expectedSequence {
		t.Errorf("Expected: %v, Got %v", expectedSequence, testApplication.sequence)
	}
}