package server

import (
	"context"
	"fmt"
	"swd_project/src/db/postgresdb"
	"swd_project/src/model"
	"swd_project/src/pbs/userpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
}

func (*UserServer) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	var user model.User
	if err := postgresdb.DB.Where("phone = ? AND password = ?", req.Phone, req.Password).Find(&user).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while fetching data from database : %v", err),
		)
	}
	if user.ID == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("شماره یا رمز عبور نادرست است"),
		)
	}
	return &userpb.LoginResponse{
		User: &userpb.User{
			Id:       int32(user.Model.ID),
			FullName: user.Name,
			Phone:    user.Phone,
		},
	}, nil
}
