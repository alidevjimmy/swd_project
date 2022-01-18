package server

import (
	"context"
	"swd_project/src/pbs/reservepb"
)

type ReserveServer struct {
	reservepb.UnimplementedReserveServiceServer
}

func (*ReserveServer) Reserve(ctx context.Context, req *reservepb.ReserveRequest) (*reservepb.ReserveResponse, error) {
	return nil, nil
}
