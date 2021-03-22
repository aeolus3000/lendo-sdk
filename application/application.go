package application

import "github.com/aeolus3000/lendo-sdk/configuration"

type Application interface {
	Initialize(configuration configuration.Configuration)
	Execute()
	Shutdown()
	setCoreApplication(coreApplication CoreApplication)
}
