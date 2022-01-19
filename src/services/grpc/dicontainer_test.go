package mygrpc_test

import (
	"testing"

	g "github.com/alikarimii/go_starter/src/services/grpc"
	"github.com/alikarimii/go_starter/src/shared"
)

func TestBuildContainer(t *testing.T) {
	l := shared.NewNilLogger()
	config := g.MustBuildConfigFromEnv(l)

	buildContainer := func() *g.DIContainer {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("[buildContainer]: %s", err)
			}
		}()
		return g.MustBuildDIContainer(config, l)
		// set options like this
		/** g.MustBuildDIContainer(
			config,
			logger,
			g.UseMongoDBConn(db),
			g.WithMarshalEvents(), ....)
		**/
	}
	// use container for other test
	_ = buildContainer()
}
