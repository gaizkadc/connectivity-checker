/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/connectivity-checker/pkg/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"time"
)

var config = server.Config{}

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
	runCmd.Flags().IntVar(&config.Port, "port", 8384, "port where connectivity-checker listens to")
	runCmd.Flags().StringVar(&config.ClusterAPIHostname, "clusterAPIHostname", "","Hostname of the cluster API on the management cluster" )
	runCmd.Flags().IntVar(&config.ClusterAPIPort, "clusterAPIPort", 8000, "Port where the cluster API is listening")
	runCmd.Flags().StringVar(&config.LoginHostname, "loginHostname", "", "Hostname of the login service")
	runCmd.Flags().IntVar(&config.LoginPort, "loginPort", 31683, "port where the login service is listening")
	runCmd.Flags().BoolVar(&config.UseTLSForLogin, "useTLSForLogin", true, "Use TLS to connect to the Login API")
	runCmd.Flags().StringVar(&config.Email, "email", "", "email address")
	runCmd.Flags().StringVar(&config.Password, "password", "", "password")
	runCmd.Flags().StringVar(&config.CACertPath, "caCertPath", "", "Path for the CA certificate")
	runCmd.Flags().StringVar(&config.ClientCertPath, "clientCertPath", "", "Path for the client certificate")
	runCmd.Flags().BoolVar(&config.SkipServerCertValidation, "skipServerCertValidation", true, "Skip CA authentication validation")
	runCmd.Flags().DurationVar(&config.ConnectivityCheckPeriod, "connectivityCheckPeriod", time.Duration(30)*time.Second, "connectivity Check Period")
	runCmd.Flags().DurationVar(&config.ConnectivityGracePeriod, "connectivityGracePeriod", time.Duration(120)*time.Second, "connectivity Grace Period")
	rootCmd.AddCommand(runCmd)
}

func RunConnectivityChecker() {
	log.Info().Msg("Launching connectivity-checker!")
	server, err := server.NewService(config)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating connectivity-checker")
	}
	server.Run()
}