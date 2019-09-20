/*
 * Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package connectivity_checker

import (
	"context"
	"github.com/nalej/grpc-infrastructure-go"
	"time"
)

func CheckCheckCheck (h *Handler, ctx context.Context, clusterId *grpc_infrastructure_go.ClusterId, duration time.Duration) {
	for true {
		h.ClusterAlive(ctx, clusterId)
		time.Sleep(duration)
	}
}
