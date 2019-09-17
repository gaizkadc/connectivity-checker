/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package server

import (
	connectivity_checker "github.com/nalej/connectivity-checker/pkg/server/connectivity-checker"
	"google.golang.org/grpc"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/reflection"
	"github.com/nalej/grpc-cluster-api-go"
	"net"
	"fmt"
)

type Service struct {
	// Server for incoming requests
	server *grpc.Server
	// Configuration object
	configuration Config
}


func NewService(config Config) (*Service, error) {
	server := grpc.NewServer()
	instance := &Service{
		server:             server,
		configuration:      config,
	}

	return instance, nil
}


func(s *Service) Run() {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.configuration.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
	}

	connectivityCheckerManager := connectivity_checker.NewManager(s.configuration)
	connectivityCheckerHandler := connectivity_checker.NewHandler(connectivityCheckerManager)

	grpcServer := grpc.NewServer()
	grpc_cluster_api_go.RegisterConnectivityCheckerServer(grpcServer, connectivityCheckerHandler)

	// Register reflection service on gRPC server
	if s.configuration.Debug {
		reflection.Register(s.server)
	}

	// Run
	log.Info().Uint32("port", s.configuration.Port).Msg("Launching gRPC server")
	if err := s.server.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}

}