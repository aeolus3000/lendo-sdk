package configuration

import "github.com/kelseyhightower/envconfig"

type EnvConfiguration struct {
}

func (ec *EnvConfiguration) Process(prefix string, config interface{}) error {
	return envconfig.Process(prefix, config)
}
