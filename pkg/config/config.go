/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"github.com/nalej/connectivity-checker/version"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-connectivity-manager-go"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

type Config struct {
	// incoming port
	Port int
	// Debugging flag
	Debug bool
	// ClusterAPIHostname with the hostname of the cluster API on the management cluster
	ClusterAPIHostname string
	// ClusterAPIPort with the port where the cluster API is listening.
	ClusterAPIPort int
	// LoginHostname with the hostname of the login API on the management cluster.
	LoginHostname string
	// LoginPort with the port where the login API is listening
	LoginPort int
	// DeploymentManagerAddress
	DeploymentManagerAddress string
	// Email to log into the management cluster.
	Email string
	// Password to log into the managment cluster.
	Password string
	// Bool to check if connections will be created securely or not
	UseTLSForLogin bool
	// Path for the certificate of the CA
	CACertPath string
	// Client Cert Path
	ClientCertPath string
	// Skip Server validation
	SkipServerCertValidation bool
	// ConnectivityCheckPeriod
	ConnectivityCheckPeriod time.Duration
	// ConnectivityGracePeriod
	ConnectivityGracePeriod time.Duration
	// Cluster ID
	ClusterId string
	// Organization ID
	OrganizationId string
	// Offline Policy must be set to true when a cluster is offline thus an offline policy should be triggered
	OfflinePolicy grpc_connectivity_manager_go.OfflinePolicy
}

func (conf *Config) Validate() derrors.Error {
	if conf.Port == 0 {
		return derrors.NewInvalidArgumentError("port must be set")
	}
	if conf.Email == "" {
		return derrors.NewInvalidArgumentError("email must be set")
	}
	if conf.Password == "" {
		return derrors.NewInvalidArgumentError("password must be set")
	}
	if conf.ClusterId == "" {
		return derrors.NewInvalidArgumentError("cluster id must be set")
	}
	if conf.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization id must be set")
	}
	if conf.ClientCertPath == "" {
		return derrors.NewInvalidArgumentError("client cert path must be set")
	}
	if conf.CACertPath == "" {
		return derrors.NewInvalidArgumentError("ca cert path must be set")
	}
	if conf.ClusterAPIHostname == "" {
		return derrors.NewInvalidArgumentError("cluster-api hostname must be set")
	}
	if conf.LoginHostname == "" {
		return derrors.NewInvalidArgumentError("login hostname must be set")
	}
	if conf.DeploymentManagerAddress == "" {
		return derrors.NewInvalidArgumentError("deployment-manager address must be set")
	}

	return nil
}

func (conf *Config) Print() {
	log.Info().Str("app", version.AppVersion).Str("commit", version.Commit).Msg("Version")
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	log.Info().Str("cluster api hostname", conf.ClusterAPIHostname).Int("port", conf.ClusterAPIPort).Msg("Cluster API on management cluster")
	log.Info().Str("login hostname", conf.LoginHostname).Int("port", conf.LoginPort).Bool("UseTLSForLogin", conf.UseTLSForLogin).Msg("Login API on management cluster")
	log.Info().Str("deployment-manager hostname", conf.DeploymentManagerAddress).Msg("Deployment Manager address")
	log.Info().Str("email", conf.Email).Str("password", strings.Repeat("*", len(conf.Password))).Msg("Application cluster credentials")
	log.Info().Str("cluster id", conf.ClusterId).Msg("cluster id")
	log.Info().Str("organization id ", conf.OrganizationId).Msg("organization id")
}
