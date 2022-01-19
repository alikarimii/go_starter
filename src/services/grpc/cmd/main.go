package main

import (
	"os"

	mygrpc "github.com/alikarimii/go_starter/src/services/grpc"
	"github.com/alikarimii/go_starter/src/shared"
)

func main() {
	stdLogger := shared.NewStandardLogger()
	config := mygrpc.MustBuildConfigFromEnv(stdLogger)
	exitFn := func() { os.Exit(1) }
	// @TODO check nil reference error if db connection be nil
	// postgresDBConn := mygrpc.MustInitPostgresDB(config, stdLogger)
	mongoDBConn := mygrpc.MustInitMongoDB(config, stdLogger)
	diContainer := mygrpc.MustBuildDIContainer(
		config,
		stdLogger,
		// mygrpc.UsePostgresDBConn(postgresDBConn), for postgres
		mygrpc.UseMongoDBConn(mongoDBConn),
	)

	s := mygrpc.InitService(config, stdLogger, exitFn, diContainer)
	go s.StartGRPCServer()
	s.WaitForStopSignal()
}
