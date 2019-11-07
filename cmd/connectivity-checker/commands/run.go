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

package commands

import (
	"github.com/nalej/connectivity-checker/pkg/config"
	"github.com/nalej/connectivity-checker/pkg/server"
	"github.com/nalej/grpc-connectivity-manager-go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

var conf = config.Config{}
var policyName string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run connectivity-checker",
	Long:  `Run connectivity-checker`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		RunConnectivityChecker()
	},
}

func init() {
	runCmd.Flags().IntVar(&conf.Port, "port", 8384, "port where connectivity-checker listens to")
	runCmd.Flags().StringVar(&conf.ClusterAPIHostname, "clusterAPIHostname", "", "Hostname of the cluster API on the management cluster")
	runCmd.Flags().IntVar(&conf.ClusterAPIPort, "clusterAPIPort", 8000, "Port where the cluster API is listening")
	runCmd.Flags().StringVar(&conf.LoginHostname, "loginHostname", "", "Hostname of the login service")
	runCmd.Flags().IntVar(&conf.LoginPort, "loginPort", 31683, "port where the login service is listening")
	runCmd.Flags().StringVar(&conf.DeploymentManagerAddress, "deploymentManagerAddress", "", "Address of the deployment-manager service")
	runCmd.Flags().BoolVar(&conf.UseTLSForLogin, "useTLSForLogin", true, "Use TLS to connect to the Login API")
	runCmd.Flags().StringVar(&conf.Email, "email", "", "email address")
	runCmd.Flags().StringVar(&conf.Password, "password", "", "password")
	runCmd.Flags().StringVar(&conf.CACertPath, "caCertPath", "", "Path for the CA certificate")
	runCmd.Flags().StringVar(&conf.ClientCertPath, "clientCertPath", "", "Path for the client certificate")
	runCmd.Flags().StringVar(&conf.ClusterId, "clusterId", "", "Cluster ID")
	runCmd.Flags().StringVar(&conf.OrganizationId, "organizationId", "", "Organization ID")
	runCmd.Flags().BoolVar(&conf.SkipServerCertValidation, "skipServerCertValidation", true, "Skip CA authentication validation")
	runCmd.Flags().DurationVar(&conf.ConnectivityCheckPeriod, "connectivityCheckPeriod", time.Duration(30)*time.Second, "connectivity Check Period")
	runCmd.Flags().DurationVar(&conf.ConnectivityGracePeriod, "connectivityGracePeriod", time.Duration(120)*time.Second, "connectivity Grace Period")
	runCmd.Flags().StringVar(&policyName, "offlinePolicy", "none", "Offline policy to trigger when cordoning an offline cluster: none or drain")
	rootCmd.AddCommand(runCmd)
}

func RunConnectivityChecker() {
	policy, exists := grpc_connectivity_manager_go.OfflinePolicy_value[strings.ToUpper(policyName)]
	if !exists {
		log.Error().Msg("invalid offline policy set")
	}
	conf.OfflinePolicy = grpc_connectivity_manager_go.OfflinePolicy(policy)

	log.Info().Msg("Launching connectivity-checker!")
	server, err := server.NewService(conf)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating connectivity-checker")
	}
	server.Run()
}
