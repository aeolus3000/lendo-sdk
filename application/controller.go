package application

import "github.com/aeolus3000/lendo-sdk/configuration"

type Controller interface {
	execute(configuration configuration.Configuration)
}
