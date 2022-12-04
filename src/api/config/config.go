package config

import (
	"dev.azure.com/SiigoDevOps/Siigo/_git/Siigo.Golang.Configuration.git/configuration"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

const KEY_ENV = "GO_ENV"
const defaultEnv = "Local"

func NewViperConfig() *viper.Viper {

	var environment string

	// read environment variable
	if env, ok := os.LookupEnv(KEY_ENV); ok {
		environment = env
	} else {
		environment = defaultEnv
	}

	return configuration.NewViperBuilder().
		WithConfigType("yaml").
		WithBasePath("configuration").
		WithConfigFile("appsettings", true).
		WithConfigFile(fmt.Sprintf("appsettings.%s", environment), true).
		WithSpringCloudConfig().
		Build()
}

func NewConfiguration(v *viper.Viper) *Configuration {

	c := &Configuration{}
	err := v.Unmarshal(c)
	if err != nil {
		panic(err)
	}

	return c
}
