package configuration

type Configuration interface {
	Process(prefix string, config interface{}) error
}

