package mygrpc

import (
	"database/sql"

	customergrpc "github.com/alikarimii/go_starter/src/customeraccounts/infrastructure/adapter/grpc"
	pb "github.com/alikarimii/go_starter/src/customeraccounts/infrastructure/adapter/grpc/proto"
	"github.com/alikarimii/go_starter/src/shared"
	"github.com/alikarimii/go_starter/src/shared/es"
	"github.com/cockroachdb/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type DIOption func(container *DIContainer) error

type DIContainer struct {
	config *Config

	infra struct {
		// db connections
		mongodbConn *mongo.Client
		pgDBConn    *sql.DB
	}

	service struct {
		// application service layer, ...
		grpcServer         *grpc.Server
		grpcCustomerServer pb.CustomerServer
	}

	dependency struct {
		// shared,marshaler, ...
		marshalEvent es.MarshalDomainEvent
	}
}

func (container *DIContainer) init() {
	container.GetEventStore()
	container.GetSomeRepository()
	container.GetCommandHandler()
	container.GetQueryHandler()
	container.GetGRPCServer()
}

// implement this
func (container *DIContainer) GetEventStore()     {}
func (container *DIContainer) GetSomeRepository() {}
func (container *DIContainer) GetCommandHandler() {}
func (container *DIContainer) GetQueryHandler()   {}
func (container *DIContainer) GetMongoDBConn() *mongo.Client {
	return container.infra.mongodbConn
}
func (container *DIContainer) GetPostgresDBConn() *sql.DB {
	return container.infra.pgDBConn
}
func (container *DIContainer) GetGRPCServer() *grpc.Server {
	if container.service.grpcServer == nil {
		container.service.grpcServer = grpc.NewServer()
		pb.RegisterCustomerServer(container.service.grpcServer, container.getGRPCCustomerServer())
		reflection.Register(container.service.grpcServer)
	}

	return container.service.grpcServer
}

// modify container base on requirements
func UseMongoDBConn(dbConn *mongo.Client) DIOption {
	return func(container *DIContainer) error {
		if dbConn == nil { // && if your db is mongo
			return errors.New("mongodbConn must not be nil")
		}

		container.infra.mongodbConn = dbConn

		return nil
	}
}

func WithMarshalEvents(fn es.MarshalDomainEvent) DIOption {
	return func(container *DIContainer) error {
		container.dependency.marshalEvent = fn
		return nil
	}
}

func MustBuildDIContainer(config *Config, l *shared.Logger, opts ...DIOption) *DIContainer {
	container := &DIContainer{}
	container.config = config

	/*** Define default dependencies ***/
	// replays marshalEvent func with actual implementation
	// or use WithMarshalEvents in test file
	container.dependency.marshalEvent = func(event es.DomainEvent) ([]byte, error) { return []byte(""), nil }

	/*** Apply options for infra, dependencies, services ***/
	for _, opt := range opts {
		if err := opt(container); err != nil {
			l.Panic().Msgf("mustBuildDIContainer: %s", err)
		}
	}

	container.init()

	return container
}

func (container *DIContainer) getGRPCCustomerServer() pb.CustomerServer {
	if container.service.grpcCustomerServer == nil {
		container.service.grpcCustomerServer = customergrpc.NewCustomerServer()
	}

	return container.service.grpcCustomerServer
}
