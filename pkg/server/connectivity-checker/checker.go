/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package connectivity_checker

import (
	"context"
	"github.com/nalej/connectivity-checker/pkg/Config"
	"github.com/nalej/connectivity-checker/pkg/login_helper"
	"github.com/nalej/grpc-cluster-api-go"
	grpc_common_go "github.com/nalej/grpc-common-go"
	grpc_connectivity_manager_go "github.com/nalej/grpc-connectivity-manager-go"
	grpc_deployment_manager_go "github.com/nalej/grpc-deployment-manager-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"
	"time"
)

func CheckClusterConnectivity (connectivityCheckerClient grpc_cluster_api_go.ConnectivityCheckerClient, clusterAPILoginHelper login_helper.LoginHelper, clusterId *grpc_infrastructure_go.ClusterId, duration time.Duration, opClient grpc_deployment_manager_go.OfflinePolicyClient,conf Config.Config) {
	var lastAliveTimestamp time.Time

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

		// Check grace period expiration
		if time.Now().Unix() - lastAliveTimestamp.Unix() > int64(conf.ConnectivityGracePeriod.Seconds()) {
			triggerOfflinePolicy(conf, opClient)
		}
		time.Sleep(duration)
	}
}

// Checks if an OfflinePolicy is set and acts accordingly
func triggerOfflinePolicy (conf Config.Config, opClient grpc_deployment_manager_go.OfflinePolicyClient) {
	log.Debug().Msg("triggering offline policy")
	switch conf.OfflinePolicy {
	case grpc_connectivity_manager_go.OfflinePolicy_NONE:
		log.Debug().Str("offline policy", conf.OfflinePolicy.String()).Msg("offline policy set to none, no additional steps required")
	case grpc_connectivity_manager_go.OfflinePolicy_DRAIN:
		remCtx, remCancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer remCancel()

		_, remErr := opClient.RemoveAll(remCtx, &grpc_common_go.Empty{})
		if remErr != nil {
			log.Error().Err(remErr).Msg("error trying to remove all instances")
		}
	default:
		log.Debug().Msg("offline policy not set, doing nothing")
	}
}
