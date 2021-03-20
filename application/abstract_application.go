package application

type AbstractApplication struct {
	Application
	coreApplication CoreApplication
}

func (aa *AbstractApplication) setCoreApplication(coreApplication CoreApplication) {
	aa.coreApplication = coreApplication
}

func (aa *AbstractApplication) GetArgs() []string {
	return aa.coreApplication.args
}