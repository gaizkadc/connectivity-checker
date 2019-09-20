/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package connectivity_checker

import (
	"github.com/nalej/connectivity-checker/pkg/login_helper"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-cluster-api-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/rs/zerolog/log"
)

type Manager struct {
	ClusterAPIClient grpc_cluster_api_go.ConnectivityCheckerClient
	ClusterAPILoginHelper *login_helper.LoginHelper
}

func NewManager(clusterAPIClient grpc_cluster_api_go.ConnectivityCheckerClient, clusterAPILoginHelper *login_helper.LoginHelper) Manager {
	return Manager{
		ClusterAPIClient: clusterAPIClient,
		ClusterAPILoginHelper:      clusterAPILoginHelper,
	}
}

func (m *Manager) ClusterAlive (clusterId *grpc_infrastructure_go.ClusterId) (*grpc_common_go.Success, derrors.Error) {
	ctx, cancel := m.ClusterAPILoginHelper.GetContext()
	defer cancel()

	result, err := m.ClusterAPIClient.ClusterAlive(ctx, clusterId)
	if err != nil {
		log.Error().Err(err).Msg("cluster doesn't seem to be alive")
	} else {
		log.Debug().Msg("cluster is alive")
	}

	return result, nil
}