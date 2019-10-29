/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package connectivity_checker

import (
	"github.com/nalej/connectivity-checker/pkg/config"
	grpc_deployment_manager_go "github.com/nalej/grpc-deployment-manager-go"
	"time"
)

const (
	DefaultTimeout =  2*time.Minute
)

// Manager structure with the remote clients required
type Manager struct {}

// NewManager creates a new manager.
func NewManager(opClient *grpc_deployment_manager_go.OfflinePolicyClient, config config.Config) (*Manager, error) {
	return &Manager{}, nil
}