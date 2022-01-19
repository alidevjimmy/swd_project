package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"swd_project/src/db/postgresdb"
	"swd_project/src/model"
	"swd_project/src/pbs/reservepb"
)

type ReserveServer struct {
	reservepb.UnimplementedReserveServiceServer
}

func (*ReserveServer) Reserve(ctx context.Context, req *reservepb.ReserveRequest) (*reservepb.ReserveResponse, error) {
	var user model.User
	if err := postgresdb.DB.Where("id = ? AND password = ?", req.GetUserId(), req.GetCurrentPassword()).Find(&user).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while fetching data from database : %v", err),
		)
	}
	if user.ID == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("کاربر یافت نشد یا رمز عبور نادرست است"),
		)
	}
	var consultant model.Consultant
	if err := postgresdb.DB.Where("id = ?", req.GetConsultantId()).Find(&consultant).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while fetching data from database : %v", err),
		)
	}
	if consultant.ID == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("مشاور یافت نشد"),
		)
	}
	//var reserve = model.Reserve{
	//	Start:        req.Start,
	//	UserID:       0,
	//	ConsultantID: 0,
	//}
	return nil, nil
}
