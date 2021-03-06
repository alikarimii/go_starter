package mygrpc

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/alikarimii/go_starter/src/shared"
)

type Service struct {
	config       *Config
	logger       *shared.Logger
	diContainter *DIContainer
	exitFn       func()
}

func InitService(
	config *Config,
	logger *shared.Logger,
	exitFn func(),
	diContainter *DIContainer,
) *Service {

	return &Service{
		config:       config,
		logger:       logger,
		exitFn:       exitFn,
		diContainter: diContainter,
	}
}

func (s *Service) StartGRPCServer() {
	s.logger.Info().Msg("configuring gRPC server ...")

	listener, err := net.Listen("tcp", s.config.GRPC.HostAndPort)
	if err != nil {
		s.logger.Error().Msgf("failed to listen: %v", err)
		s.shutdown()
	}

	s.logger.Info().Msgf("starting gRPC server listening at %s ...", s.config.GRPC.HostAndPort)

	grpcServer := s.diContainter.GetGRPCServer()
	if err := grpcServer.Serve(listener); err != nil {
		s.logger.Error().Msgf("gRPC server failed to serve: %s", err)
		s.shutdown()
	}
}

func (s *Service) WaitForStopSignal() {
	s.logger.Info().Msg("start waiting for stop signal ...")

	stopSignalChannel := make(chan os.Signal, 1)
	signal.Notify(stopSignalChannel, os.Interrupt, syscall.SIGTERM)

	sig := <-stopSignalChannel

	s.logger.Info().Msgf("received '%s'", sig)
	close(stopSignalChannel)
	s.shutdown()
}

func (s *Service) shutdown() {
	s.logger.Info().Msg("shutdown: stopping services ...")

	grpcServer := s.diContainter.GetGRPCServer()
	if grpcServer != nil {
		s.logger.Info().Msg("shutdown: stopping gRPC server gracefully ...")
		grpcServer.GracefulStop()
	}
	{ // if used postgres
		postgresDBConn := s.diContainter.GetPostgresDBConn()
		if postgresDBConn != nil {
			s.logger.Info().Msg("shutdown: closing Postgres DB connection ...")
			if err := postgresDBConn.Close(); err != nil {
				s.logger.Warn().Msgf("shutdown: failed to close the Postgres DB connection: %s", err)
			}
		}
	}

	{ // if used mongodb
		mongodbConn := s.diContainter.GetMongoDBConn()
		if mongodbConn != nil {
			s.logger.Info().Msg("shutdown: closing Mongodb DB connection ...")
			if err := mongodbConn.Disconnect(context.Background()); err != nil {
				s.logger.Warn().Msgf("shutdown: failed to close the Mongodb DB connection: %s", err)
			}
		}
	}

	s.logger.Info().Msg("shutdown: all services stopped!")

	s.exitFn()
}
