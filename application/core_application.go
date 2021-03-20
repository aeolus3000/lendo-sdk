package application

import "lendo-sdk/configuration"

type CoreApplication struct {
	configuration configuration.Configuration
	name string
	args []string
}

func NewCoreApplication(name string, args []string) CoreApplication {
	return CoreApplication{
		configuration: configuration.NewDefaultConfiguration(),
		name:          name,
		args:          args,
	}
}

func (ca *CoreApplication) getConfiguration() configuration.Configuration {
	return ca.configuration
}