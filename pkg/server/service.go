/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/nalej/connectivity-checker/pkg/config"
	"github.com/nalej/connectivity-checker/pkg/login_helper"
	"github.com/nalej/connectivity-checker/pkg/server/connectivity-checker"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-cluster-api-go"
	"github.com/nalej/grpc-deployment-manager-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-login-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"net"
)

type Service struct {
	// Server for incoming requests
	Server *grpc.Server
	// Configuration object
	Configuration config.Config
}

func NewService(config config.Config) (*Service, error) {
	server := grpc.NewServer()
	service := &Service{
		Server:             server,
		Configuration:      config,
	}

	return service, nil
}

type Clients struct {
	ConnectivityCheckerClient grpc_cluster_api_go.ConnectivityCheckerClient
	LoginClient  grpc_login_api_go.LoginClient
	OfflinePolicyClient grpc_deployment_manager_go.OfflinePolicyClient
}

func (s *Service) GetClients() (*Clients, derrors.Error) {

	ccConn, err := s.getSecureAPIConnection(s.Configuration.ClusterAPIHostname, s.Configuration.ClusterAPIPort, s.Configuration.CACertPath, s.Configuration.ClientCertPath, s.Configuration.SkipServerCertValidation)
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the Cluster API manager")
	}
	connectivityCheckerClient := grpc_cluster_api_go.NewConnectivityCheckerClient(ccConn)

	loginConn, err := s.getSecureAPIConnection(s.Configuration.LoginHostname, s.Configuration.LoginPort, s.Configuration.CACertPath, s.Configuration.ClientCertPath, s.Configuration.SkipServerCertValidation)
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the Login API manager")
	}
	loginClient := grpc_login_api_go.NewLoginClient(loginConn)

	opConn, opErr := grpc.Dial(s.Configuration.DeploymentManagerAddress, grpc.WithInsecure())
	if opErr != nil {
		return nil, derrors.AsError(opErr, "cannot create connection with the Deployment Manager")
	}
	opClient := grpc_deployment_manager_go.NewOfflinePolicyClient(opConn)

	return &Clients{ConnectivityCheckerClient:connectivityCheckerClient, LoginClient:loginClient, OfflinePolicyClient:opClient}, nil
}

func (s *Service) getSecureAPIConnection(hostname string, port int, caCertPath string, clientCertPath string, skipCAValidation bool) (*grpc.ClientConn, derrors.Error) {
	// Build connection with cluster API
	rootCAs := x509.NewCertPool()
	tlsConfig := &tls.Config{
		ServerName:   hostname,
	}

	if caCertPath != "" {
		log.Debug().Str("caCertPath", caCertPath).Msg("loading CA cert")
		caCert, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			return nil, derrors.NewInternalError("Error loading CA certificate")
		}
		added := rootCAs.AppendCertsFromPEM(caCert)
		if !added {
			return nil, derrors.NewInternalError("cannot add CA certificate to the pool")
		}
		tlsConfig.RootCAs = rootCAs
	}

	targetAddress := fmt.Sprintf("%s:%d", hostname, port)
	log.Debug().Str("address", targetAddress).Msg("creating cluster API connection")

	if clientCertPath != "" {
		log.Debug().Str("clientCertPath", clientCertPath).Msg("loading client certificate")
		clientCert, err := tls.LoadX509KeyPair(fmt.Sprintf("%s/tls.crt", clientCertPath),fmt.Sprintf("%s/tls.key", clientCertPath))
		if err != nil {
			log.Error().Str("error", err.Error()).Msg("Error loading client certificate")
			return nil, derrors.NewInternalError("Error loading client certificate")
		}

		tlsConfig.Certificates = []tls.Certificate{clientCert}
		tlsConfig.BuildNameToCertificate()
	}
	log.Debug().Str("address", targetAddress).Str("caCertPath", caCertPath).Bool("skipCAValidation", skipCAValidation).Msg("creating secure connection")

	if skipCAValidation {
		log.Debug().Msg("skipping server cert validation")
		tlsConfig.InsecureSkipVerify = true
	}

	creds := credentials.NewTLS(tlsConfig)

	log.Debug().Interface("creds", creds.Info()).Msg("Secure credentials")
	sConn, dErr := grpc.Dial(targetAddress, grpc.WithTransportCredentials(creds))
	if dErr != nil {
		return nil, derrors.AsError(dErr, "cannot create connection with the cluster API service")
	}
	return sConn, nil
}

func(s *Service) Run () {
	cErr := s.Configuration.Validate()
	if cErr != nil {
		log.Fatal().Str("err", cErr.DebugReport()).Msg("invalid configuration")
	}
	s.Configuration.Print()

	// create clients
	clients, cErr := s.GetClients()
	if cErr != nil {
		log.Fatal().Str("err", cErr.DebugReport()).Msg("cannot generate clients")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Configuration.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
	}

	connectivityCheckerClient := clients.ConnectivityCheckerClient
	clusterAPILoginHelper := login_helper.NewLogin(s.Configuration.LoginHostname, s.Configuration.LoginPort, s.Configuration.UseTLSForLogin,
		s.Configuration.Email, s.Configuration.Password, s.Configuration.CACertPath, s.Configuration.ClientCertPath, s.Configuration.SkipServerCertValidation)
	opClient := clients.OfflinePolicyClient

	lErr := clusterAPILoginHelper.Login()
	if lErr != nil {
		log.Fatal().Str("err", cErr.DebugReport()).Msg("there was an error requesting cluster-api login")
	}

	// Register reflection service on gRPC server
	if s.Configuration.Debug {
		reflection.Register(s.Server)
	}

	// Infinite loop of ClusterAlive signalsa and grace expiration checks
	log.Debug().Str("cluster id", s.Configuration.ClusterId).Msg("cluster id")
	log.Debug().Dur("connectivity check period", s.Configuration.ConnectivityCheckPeriod).Msg("ConnectivityCheckPeriod")
	clusterId :=  &grpc_infrastructure_go.ClusterId{
		ClusterId: s.Configuration.ClusterId,
		OrganizationId: s.Configuration.OrganizationId,
	}
	go connectivity_checker.CheckClusterConnectivity(connectivityCheckerClient, *clusterAPILoginHelper, clusterId, s.Configuration.ConnectivityCheckPeriod, opClient, s.Configuration)

	// Run
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	err = s.Server.Serve(lis)
	if err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
}