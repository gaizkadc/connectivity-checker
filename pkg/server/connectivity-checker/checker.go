/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package connectivity_checker

import (
	"github.com/nalej/connectivity-checker/pkg/login_helper"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-cluster-api-go"
	"github.com/rs/zerolog/log"
	"time"
)

func CheckClusterConnectivity (connectivityCheckerClient grpc_cluster_api_go.ConnectivityCheckerClient, clusterAPILoginHelper login_helper.LoginHelper, clusterId *grpc_infrastructure_go.ClusterId, duration time.Duration) {
	for true {
		ctx, _ := clusterAPILoginHelper.GetContext()
		_, err := connectivityCheckerClient.ClusterAlive(ctx, clusterId)
		if err != nil {
			log.Error().Err(err).Msg("cluster doesn't seem to be alive")
		} else {
			log.Debug().Msg("cluster is alive")
		}
		time.Sleep(duration)
	}
}
