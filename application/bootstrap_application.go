package application

type BootstrapApplication struct {
	application Application
	name        string
}

func NewBootstrapApplication(application Application, name string) BootstrapApplication {
	return BootstrapApplication{application, name}
}

func (b BootstrapApplication) Execute(args []string) {
	coreApplication := NewCoreApplication(b.name, args)
	b.application.setCoreApplication(coreApplication)

	controller := NewDefaultController(b.application)
	controller.execute(coreApplication.configuration)
}

func (b BootstrapApplication) WithApplicationName(name string) {
	b.name = name
}
