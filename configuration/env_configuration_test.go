package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

func ExampleEnvConfiguration_Process() {
	type TestConfigInner struct {
		ManualOverride1 string `envconfig:"manual_override_1"` //looks for MANUAL_OVERRIDE_1
		DefaultVar      string `default:"foobar"` //if no ENV found, the value provided will be used
		RequiredVar     string `required:"true"` //mandatory field
		IgnoredVar      string `ignored:"true"` //will be ignored even if an ENV exists
		AutoSplitVar    string `split_words:"true"` //looks for AUTO_SPLIT_VAR
		RequiredAndAutoSplitVar    string `required:"true" split_words:"true"` //combination of properties
	}
	type TestConfig struct {
		Value string `default:"foobar"`
		Another TestConfigInner
	}
	testConfig := TestConfig{}
	conf := NewDefaultConfiguration()
	prefix := "PREFIX"
	setEnv(prefix + "_VALUE", "helloworld")
	setEnv(prefix + "_ANOTHER_VALUEINNER", "worldhello")
	setEnv(prefix + "_ANOTHER_REQUIREDVAR", "hweolrldo")
	setEnv(prefix + "_ANOTHER_REQUIRED_AND_AUTO_SPLIT_VAR", "abcdefgh")
	err := conf.Process(prefix, &testConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonBytes, err := json.Marshal(testConfig)
	fmt.Println(string(jsonBytes))
	//Output:
	//{"Value":"helloworld","Another":{"ManualOverride1":"","DefaultVar":"foobar","RequiredVar":"hweolrldo","IgnoredVar":"","AutoSplitVar":"","RequiredAndAutoSplitVar":"abcdefgh"}}
}

func setEnv(key, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		fmt.Println(err.Error())
	}
}