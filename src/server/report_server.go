package server

import (
	"context"
	"fmt"
	"swd_project/src/db/postgresdb"
	"swd_project/src/model"
	"swd_project/src/pbs/reportpb"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ReportServer struct {
	reportpb.UnimplementedReportServiceServer
}

func (*ReportServer) CreateReport(ctx context.Context, req *reportpb.CreateReportRequest) (*reportpb.CreateReportResponse, error) {
	var user model.User
	if err := postgresdb.DB.Where("id = ? AND password = ?", req.GetReport().GetUserId(), req.GetCurrentPassword()).Find(&user).Error; err != nil {
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
	report := model.Report{
		Subject: req.GetReport().GetSubject(),
		Long:    int(req.GetReport().GetLong()),
		Lat:     int(req.GetReport().GetLat()),
		Address: req.GetReport().GetAddress(),
		Active:  true,
		UserID:  int(req.GetReport().GetUserId()),
		Until:   req.GetReport().GetUntil().AsTime(),
	}
	if err := postgresdb.DB.Create(&report).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("خطا هنگام ایجاد گزارش"),
		)
	}
	if req.GetReport().GetUntil() == nil {
		user.Status = model.Red
	} else {
		if report.Until.Unix() < time.Now().Unix() {
			return nil, status.Errorf(
				codes.Internal,
				fmt.Sprintf("زمان احتمال می دهید کی در خطر باشید باید از الان بیشتر باشد"),
			)
		}
		user.Status = model.Yellow
	}
	if err := postgresdb.DB.Save(&user).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("خطا هنگاه بروزرسانی وضعیت کاربر"),
		)
	}
	until := timestamppb.New(report.Until)
	return &reportpb.CreateReportResponse{
		Report: &reportpb.Report{
			Id:      int32(report.ID),
			Subject: report.Subject,
			Long:    int32(report.Long),
			Lat:     int32(report.Lat),
			Address: report.Address,
			UserId:  int32(report.UserID),
			Until:   until,
		},
	}, nil
}

func (*ReportServer) UserOpenReports(ctx context.Context, req *reportpb.UserOpenReportsRequest) (*reportpb.UserOpenReportsResponse, error) {
	var user model.User
	if err := postgresdb.DB.Where("id = ?", req.GetUserId()).Find(&user).Error; err != nil {
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
	var reports []model.Report
	if err := postgresdb.DB.Where("active = ?", true).Find(&reports).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while fetching data from database : %v", err),
		)
	}
	var reportsRes = []*reportpb.Report{}
	for _, v := range reports {
		until := timestamppb.New(v.Until)
		reportsRes = append(reportsRes, &reportpb.Report{
			Id:      int32(v.Model.ID),
			Subject: v.Subject,
			Long:    int32(v.Long),
			Lat:     int32(v.Lat),
			Address: v.Address,
			UserId:  int32(v.UserID),
			Until:   until,
		})
	}
	return &reportpb.UserOpenReportsResponse{
		Reports: reportsRes,
	}, nil
}
