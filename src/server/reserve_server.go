package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"swd_project/src/db/postgresdb"
	"swd_project/src/model"
	"swd_project/src/pbs/reservepb"
	"time"
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
	var userReserves []model.Reserve
	if err := postgresdb.DB.Where("consultant_id = ? AND start > ? AND start < ?", req.GetConsultantId(), req.GetStart().AsTime().Add(-24*time.Hour), req.GetStart().AsTime().Add(24*time.Hour)).Find(&userReserves).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while fetching data from database : %v", err),
		)
	}
	if len(userReserves) != 0 {
		return nil, status.Errorf(
			codes.FailedPrecondition,
			fmt.Sprintf("در ۲۴ ساعت فقط یک وقت مشاوره می توانید رزرو کنید"),
		)
	}
	var reserve = model.Reserve{
		Start:        req.GetStart().AsTime(),
		UserID:       int(req.GetUserId()),
		ConsultantID: int(req.GetConsultantId()),
	}
	if err := postgresdb.DB.Create(&reserve).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("خطا هنگام ایجاد رزرو"),
		)
	}
	return &reservepb.ReserveResponse{}, nil
}
