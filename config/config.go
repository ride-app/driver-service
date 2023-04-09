package config

type ConfigStruct struct {
	Debug               bool   `env:"DEBUG" env-description:"dev or prod" env-default:"false"`
	Port                int32  `env:"PORT" env-description:"server port" env-default:"50051"`
	Wallet_Service_Host string `env:"WALLET_SERVICE_HOST" env-description:"wallet service host" env-default:"localhost:50052"`
	Firebase_Project_Id string `env:"FIREBASE_PROJECT_ID" env-description:"firebase project id" env-default:"NO_PROJECT"`
}

var Env ConfigStruct
