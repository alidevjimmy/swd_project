package server

import (
	"context"
	"fmt"
	"swd_project/src/db/postgresdb"
	"swd_project/src/model"
	"swd_project/src/pbs/schudulepb"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SchuduleServer struct {
	schudulepb.UnimplementedSchuduleServiceServer
}

func (*SchuduleServer) CreateSchudule(ctx context.Context, req *schudulepb.CreateSchuduleRequest) (*schudulepb.CreateSchuduleResponse, error) {
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
	var schudules = []model.Schudule{}
	for _, v := range req.Schudules {
		schudules = append(schudules, model.Schudule{
			Date:         v.GetDate().AsTime(),
			StartHour:    int(v.GetStartHour()),
			FinishHour:   int(v.GetFinishHour()),
			ConsultantID: int(req.GetConsultantId()),
		})
	}
	if err := postgresdb.DB.Create(&schudules).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("خطا هنگام ذخیره سازی وقت ها"),
		)
	}
	var schudulesRes = []*schudulepb.Schudule{}
	for _, v := range schudules {
		date := timestamppb.New(v.Date)
		schudulesRes = append(schudulesRes, &schudulepb.Schudule{
			Id:           int32(v.ID),
			Date:         date,
			StartHour:    int32(v.StartHour),
			FinishHour:   int32(v.FinishHour),
			ConsultantId: int32(v.ConsultantID),
		})
	}
	return &schudulepb.CreateSchuduleResponse{
		Schudules: schudulesRes,
	}, nil
}

func (*SchuduleServer) FindAllOpenSchudules(ctx context.Context, req *schudulepb.FindAllOpenSchudulesRequest) (*schudulepb.FindAllOpenSchudulesResponse, error) {
	// something missed here (reserved times)!
	var schudules = []model.Schudule{}
	if err := postgresdb.DB.Where("date >= ?", time.Now()).Find(&schudules).Error; err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("خطا هنگام استخراج وقت های آزاد"),
		)
	}
	var schudulesRes = []*schudulepb.Schudule{}
	for _, v := range schudules {
		if v.Date == time.Now() {
			if time.Now().Local().Hour() >= v.FinishHour {
				continue
			}
		}
		date := timestamppb.New(v.Date)
		schudulesRes = append(schudulesRes, &schudulepb.Schudule{
			Id:           int32(v.ID),
			Date:         date,
			StartHour:    int32(v.StartHour),
			FinishHour:   int32(v.FinishHour),
			ConsultantId: int32(v.ConsultantID),
		})
	}
	return &schudulepb.FindAllOpenSchudulesResponse{
		Schudules: schudulesRes,
	}, nil
}

func (*SchuduleServer) FindConsultantOpenSchudules(ctx context.Context, req *schudulepb.FindConsultantOpenSchudulesRequest) (*schudulepb.FindConsultantOpenSchudulesResponse, error) {
	var schudules = []model.Schudule{}
	if err := postgresdb.DB.Where("consultant_id = ? AND date >= ?", req.GetConsultantId(), time.Now()).Find(&schudules).Error; err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("خطا هنگام استخراج وقت های آزاد"),
		)
	}
	var schudulesRes = []*schudulepb.Schudule{}
	for _, v := range schudules {
		if v.Date == time.Now() {
			if time.Now().Local().Hour() >= v.FinishHour {
				continue
			}
		}
		date := timestamppb.New(v.Date)
		schudulesRes = append(schudulesRes, &schudulepb.Schudule{
			Id:           int32(v.ID),
			Date:         date,
			StartHour:    int32(v.StartHour),
			FinishHour:   int32(v.FinishHour),
			ConsultantId: int32(v.ConsultantID),
		})
	}
	return &schudulepb.FindConsultantOpenSchudulesResponse{
		Schudules: schudulesRes,
	}, nil
}
