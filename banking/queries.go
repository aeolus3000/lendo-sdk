package banking

type BankingApi interface {
	//This method creates an application with a bank
	//In Input only firstname and lastname are mandatory
	//In Output all fields are set
	//Returns error when it can not reach out to the banking service
	Create(Application) (Application, error)

	//This method checks an application status with a bank
	//In Input only id are mandatory
	//In Output all fields are set
	//Returns error when it can not reach out to the banking service
	CheckStatus(string) (Application, error)
}
