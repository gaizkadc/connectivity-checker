/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package connectivity_checker

import (
	"github.com/nalej/connectivity-checker/pkg/login_helper"
	"github.com/nalej/grpc-cluster-api-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"
	"time"
)

func CheckClusterConnectivity (connectivityCheckerClient grpc_cluster_api_go.ConnectivityCheckerClient, clusterAPILoginHelper login_helper.LoginHelper, clusterId *grpc_infrastructure_go.ClusterId, duration time.Duration, lastAliveTimestamp time.Time) {
	for true {
		ctx, cancel := clusterAPILoginHelper.GetContext()
		if cancel != nil {
			defer cancel()
		}
		_, err := connectivityCheckerClient.ClusterAlive(ctx, clusterId)
		if err != nil {
			st := grpc_status.Convert(err).Code()
			if st == codes.Unauthenticated {
				errLogin := clusterAPILoginHelper.RerunAuthentication()
				if errLogin != nil {
					log.Error().Err(errLogin).Msg("error during reauthentication")
				}
				ctx2, cancel2 := clusterAPILoginHelper.GetContext()
				defer cancel2()
				_, err = connectivityCheckerClient.ClusterAlive(ctx2, clusterId)
				if err != nil {
					log.Error().Err(err).Msg("cluster doesn't seem to be alive")
				} else {
					log.Debug().Msg("cluster is alive")
					lastAliveTimestamp = time.Now()
				}
			}
		} else {
			log.Debug().Msg("cluster is alive")
			lastAliveTimestamp = time.Now()
		}
		time.Sleep(duration)
	}
}

