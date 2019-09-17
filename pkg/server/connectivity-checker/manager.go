/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package connectivity_checker

import (
	"github.com/nalej/connectivity-checker/pkg/server"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-infrastructure-go"
	"sync"
)

type Manager struct {
	sync.Mutex
	Config      server.Config
}

func NewManager(config server.Config) Manager {
	return Manager{
		Config:      config,
	}
}

func (m *Manager) CluserAlive (clusterId *grpc_infrastructure_go.ClusterId) (*grpc_common_go.Success, derrors.Error) {

}