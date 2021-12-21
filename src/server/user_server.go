package server

import (
	"context"
	"esme_team/src/db/postgresdb"
	"esme_team/src/model"
	"esme_team/src/pbs/userpb"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
}

func (*UserServer) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	var user model.UserModel
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
			Id:       int32(user.ID),
			FullName: user.FullName,
			Phone:    user.Phone,
		},
	}, nil
}


func (*UserServer) StoreNewToken(ctx context.Context, req *userpb.StoreNewTokenRequest) (*userpb.StoreNewTokenResponse, error) {
	model := &model.TokenModel{
		Token:  req.GetToken(),
		UserID: uint(req.GetUserId()),
	}
	if err := postgresdb.DB.Create(&model).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while storing data to database : %v", err),
		)
	}
	return &userpb.StoreNewTokenResponse{
		Result: "توکن با موفقیت افزوده شد",
	}, nil
}
