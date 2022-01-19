package mygrpc

import (
	"fmt"
	"os"

	"github.com/alikarimii/go_starter/src/shared"
)

// edit this file base on requirements of project

type Config struct {
	Mongodb struct {
		DSN                 string
		MongoInitdbDatabase string
	}
	Postgres struct {
		DSN                    string
		MigrationsPathCustomer string
	}
	GRPC struct {
		HostAndPort string
	}
	HTTP struct {
		HostAndPort string
	}
}

// ConfigExpectedEnvKeys - This is also used by Config_test.go
// to check that all keys exist in Env,
// so always add new keys here!
var ConfigExpectedEnvKeys = map[string]string{
	"postgresDSN":                    "POSTGRES_DSN",
	"mongodbDSN":                     "MONGODB_DSN",
	"postgresMigrationsPathCustomer": "POSTGRES_MIGRATIONS_PATH_CUSTOMER",
	"grpcHostAndPort":                "GRPC_HOST_AND_PORT",
	"mongoInitdbDatabase":            "MONGO_INITDB_DATABASE",
}

// this used for testing
func MustBuildConfigFromEnv(l *shared.Logger) *Config {
	var err error
	conf := &Config{}
	msg := "mustBuildConfigFromEnv: %s !"

	{ // if use mongodb as database
		if conf.Mongodb.DSN, err = conf.stringFromEnv(ConfigExpectedEnvKeys["mongodbDSN"]); err != nil {
			l.Panic().Msgf(msg, err)
		}
		if conf.Mongodb.MongoInitdbDatabase, err = conf.stringFromEnv(ConfigExpectedEnvKeys["mongoInitdbDatabase"]); err != nil {
			l.Panic().Msgf(msg, err)
		}
	}

	{ // if use postgres as database
		if conf.Postgres.DSN, err = conf.stringFromEnv(ConfigExpectedEnvKeys["postgresDSN"]); err != nil {
			l.Panic().Msgf(msg, err)
		}
		if conf.Postgres.MigrationsPathCustomer, err = conf.stringFromEnv(ConfigExpectedEnvKeys["postgresMigrationsPathCustomer"]); err != nil {
			l.Panic().Msgf(msg, err)
		}
	}

	if conf.GRPC.HostAndPort, err = conf.stringFromEnv(ConfigExpectedEnvKeys["grpcHostAndPort"]); err != nil {
		l.Panic().Msgf(msg, err)
	}

	return conf
}

func (conf Config) stringFromEnv(envKey string) (string, error) {
	envVal, ok := os.LookupEnv(envKey)
	if !ok {
		return "", fmt.Errorf("config value [%s] missing in env", envKey)
	}

	return envVal, nil
}
