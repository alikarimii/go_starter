package mygrpc_test

import (
	"testing"

	g "github.com/alikarimii/go_starter/src/services/grpc"
	"github.com/alikarimii/go_starter/src/shared"
)

func TestConfig(t *testing.T) {
	logger := shared.NewNilLogger()

	buildConfig := func() {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("[buildConfig]: %s", err)
			}
		}()
		g.MustBuildConfigFromEnv(logger)
	}
	buildConfig()
}
