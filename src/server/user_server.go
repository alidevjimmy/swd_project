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

func (*UserServer) FindUser(ctx context.Context, req *userpb.FindUserRequest) (*userpb.FindUserResponse, error) {
	var user model.User
	if err := postgresdb.DB.Where("id = ?", req.GetId()).Find(&user).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while fetching data from database : %v", err),
		)
	}
	if user.ID == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("کاربر یافت نشد"),
		)
	}
	birth := timestamppb.New(user.Birth)
	return &userpb.FindUserResponse{
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

func (*UserServer) EditUser(ctx context.Context, req *userpb.EditUserRequest) (*userpb.EditUserResponse, error) {
	var user model.User
	if err := postgresdb.DB.Where("id = ? AND password = ?", req.User.GetId(), req.GetCurrentPassword()).Find(&user).Error; err != nil {
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

	user.Password = req.User.GetPassword()
	user.Name = req.User.GetName()
	user.Family = req.User.GetFamily()
	user.NationalCode = req.User.GetNationalCode()
	user.Password = req.User.GetPassword()
	user.Birth = req.User.GetBirth().AsTime()

	if err := postgresdb.DB.Model(&user).Updates(user).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("خطا هنگاه بروزرسانی اطلاعات کاربر"),
		)
	}
	birth := timestamppb.New(user.Birth)
	return &userpb.EditUserResponse{
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

func (*UserServer) SwapStatus(ctx context.Context, req *userpb.SwapStatusRequest) (*userpb.SwapStatusResponse, error) {
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
	if user.Status != model.UserStatus(userpb.UserStatus_RED) {
		return nil, status.Errorf(
			codes.PermissionDenied,
			fmt.Sprintf("وضعیت تنها از قرمز به سبز قابل تغییر است"),
		)
	}
	user.Status = model.Green
	if err := postgresdb.DB.Save(&user).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("خطا هنگاه بروزرسانی وضعیت کاربر"),
		)
	}
	// set active field for all reports as FALSE
	return nil, nil
}
