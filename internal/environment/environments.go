package environment

import (
	"context"
	appCtx "ksef/cmd/ksef/context"
	"ksef/internal/utils"
)

type Environment string

type Config struct {
	Host         string
	Environment  Environment
	NIPValidator utils.NIPValidatorType
}

const (
	Test       Environment = "ksef-test.mf.gov.pl"
	Production Environment = "ksef.mf.gov.pl"
)

var envConfigs = map[Environment]Config{
	Test: Config{
		// feels sketchy to do it but on the other hand it saves us the string casting in
		// other parts of the code :|
		Environment:  Test,
		Host:         string(Test),
		NIPValidator: utils.NIPLengthValidator,
	},
	Production: Config{
		Environment:  Production,
		Host:         string(Production),
		NIPValidator: utils.NIPValidator,
	},
}

func GetConfig(env Environment) Config {
	return envConfigs[env]
}

func FromContext(ctx context.Context) Environment {
	envI := ctx.Value(appCtx.KeyEnvironment)
	env, ok := envI.(Environment)
	if !ok {
		panic("environment not found")
	}
	return env
}
