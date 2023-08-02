package config

type ConfigStruct struct {
	Debug             bool   `env:"DEBUG" env-description:"dev or prod" env-default:"false"`
	Port              int32  `env:"PORT" env-description:"server port" env-default:"50051"`
	ProjectID         string `env:"PROJECT_ID" env-description:"project id" env-default:"NO_PROJECT"`
	WalletServiceHost string `env:"WALLET_SERVICE_HOST" env-description:"wallet service host" env-default:"localhost:50052"`
}

var Env ConfigStruct
