package server

import (
	"context"
	"fmt"
	"strings"
	"swd_project/src/db/postgresdb"
	"swd_project/src/model"
	"swd_project/src/pbs/userpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	birth := timestamppb.New(user.Birth)
	return &userpb.LoginResponse{
		User: &userpb.User{
			Id:           int32(user.Model.ID),
			Name:         user.Name,
			Family:       user.Family,
			Phone:        user.Phone,
			NationalCode: user.NationalCode,
			UserStatus:   userpb.UserStatus(user.Status),
			Birth:        birth,
		},
	}, nil
}

func (*UserServer) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	user := model.User{
		Phone:        req.GetPhone(),
		Name:         req.GetName(),
		Family:       req.GetFamily(),
		NationalCode: req.GetNationalCode(),
		Status:       1,
		Birth:        req.GetBirth().AsTime(),
		Password:     req.GetPassword(),
	}

	if res := postgresdb.DB.Create(&user); res.Error != nil {
		if strings.Contains(res.Error.Error(), "duplicate") {
			return nil, status.Errorf(
				codes.AlreadyExists,
				fmt.Sprintf("این شماره تلفن قبلا ثبت شده است"),
			)
		}
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("خطایی رخ داده است، لطفا مجددا تلاش کنید"),
		)
	}

	birth := timestamppb.New(user.Birth)
	return &userpb.RegisterResponse{
		User: &userpb.User{
			Id:           int32(user.Model.ID),
			Phone:        req.GetPhone(),
			Name:         req.GetName(),
			Family:       req.GetFamily(),
			NationalCode: req.GetNationalCode(),
			UserStatus:   1,
			Birth:        birth,
			Password:     req.GetPassword(),
		},
	}, nil
}
