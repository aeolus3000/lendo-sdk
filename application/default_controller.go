package application

import "github.com/aeolus3000/lendo-sdk/configuration"

type DefaultController struct {
	mainApplication Application
}

func NewDefaultController(application Application) Controller {
	return &DefaultController{application}
}

func (dc *DefaultController) execute(configuration configuration.Configuration) {
	dc.mainApplication.Initialize(configuration)
	dc.mainApplication.Execute()
	dc.mainApplication.Shutdown()
}
