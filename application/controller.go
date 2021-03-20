package application

import "lendo-sdk/configuration"

type Controller interface {
	execute(configuration configuration.Configuration)
}