/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package connectivity_checker

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/rs/zerolog/log"
	"github.com/nalej/grpc-utils/pkg/conversions"
)

type Handler struct {
	Manager Manager
}

func (h Handler) ClusterAlive(ctx context.Context, clusterId *grpc_infrastructure_go.ClusterId) (*grpc_common_go.Success, error) {
	log.Debug().Str("clusterId", clusterId.ClusterId).Str("organizationId", clusterId.OrganizationId).Msg("cluster alive")
	err := h.ValidClusterId(clusterId)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	result, err := h.Manager.ClusterAlive(*clusterId)
}

func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}

func (h * Handler) ValidClusterId (clusterId *grpc_infrastructure_go.ClusterId) derrors.Error {
	if clusterId.ClusterId == "" {
		return derrors.NewInvalidArgumentError("expecting ClusterId")
	}

	if clusterId.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("expecting OrganizationId")
	}

	return nil
}