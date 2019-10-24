/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package connectivity_checker

import (
	"context"
	"github.com/nalej/connectivity-checker/pkg/Config"
	grpc_common_go "github.com/nalej/grpc-common-go"
	grpc_deployment_manager_go "github.com/nalej/grpc-deployment-manager-go"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	DefaultTimeout =  2*time.Minute
)

// Manager structure with the remote clients required
type Manager struct {
	OfflinePolicyClient grpc_deployment_manager_go.OfflinePolicyClient
	Config              Config.Config
}

// NewManager creates a new manager.
func NewManager(opClient *grpc_deployment_manager_go.OfflinePolicyClient, config Config.Config) (*Manager, error) {
	return &Manager{
		OfflinePolicyClient:*opClient,
		Config: config,
	}, nil
}

func (m *Manager) CheckGracePeriodExpiration (lastAliveTimestamp time.Time, opClient grpc_deployment_manager_go.OfflinePolicyClient, duration time.Duration) {
	for true {
		if time.Now().Unix() - lastAliveTimestamp.Unix() > int64(m.Config.ConnectivityGracePeriod.Seconds()) {
			log.Debug().Int64("now", time.Now().Unix()).Msg("time now")
			log.Debug().Int64("last alive timestamp", lastAliveTimestamp.Unix()).Msg("last alive timestamp")
			log.Debug().Int64("grace period", int64(m.Config.ConnectivityGracePeriod.Seconds())).Msg("grace period")

			remCtx, remCancel := context.WithTimeout(context.Background(), DefaultTimeout)
			defer remCancel()

			_, remErr := opClient.RemoveAll(remCtx, &grpc_common_go.Empty{})
			if remErr != nil {
				log.Error().Err(remErr).Msg("error trying to remove all instances")
			}
		}
		time.Sleep(duration)
	}
}