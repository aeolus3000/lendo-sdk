package configuration

func NewDefaultConfiguration() Configuration {
	return &EnvConfiguration{}
}