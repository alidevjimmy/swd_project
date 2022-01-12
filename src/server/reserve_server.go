package server

import (
	"context"
	"swd_project/src/pbs/reservepb"
)

type ReserveServer struct {
	reservepb.UnimplementedReserveServiceServer
}

func (*ReserveServer) ReserveOneHour(ctx context.Context, req *reservepb.ReserveOneHourRequest) (*reservepb.ReserveOneHourResponse, error) {
	// check that hour reserved or not
	// query schuduleID and check start hour
	return nil, nil
}

func (*ReserveServer) FindAllUserReserves(ctx context.Context, req *reservepb.FindAllUserReservesRequest) (*reservepb.FindAllUserReservesResponse, error) {
	//
	return nil, nil
}
