package config

type ConfigStruct struct {
	Debug bool  `env:"DEBUG" env-description:"dev or prod" env-default:"false"`
	Port  int32 `env:"PORT" env-description:"server port" env-default:"50051"`
}

var Config ConfigStruct
