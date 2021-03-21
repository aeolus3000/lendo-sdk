package application

import "lendo-sdk/configuration"

type Application interface {
	Initialize(configuration configuration.Configuration)
	Execute()
	Shutdown()
	setCoreApplication(coreApplication CoreApplication)
}
