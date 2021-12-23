package server

import (
	"context"
	"fmt"
	"swd_project/src/db/postgresdb"
	"swd_project/src/model"
	"swd_project/src/pbs/schudulepb"

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
			Id: int32(v.ID),
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
	return nil, nil
}

func (*SchuduleServer) FindConsultantOpenSchudules(ctx context.Context, req *schudulepb.FindConsultantOpenSchudulesRequest) (*schudulepb.FindConsultantOpenSchudulesResponse, error) {
	return nil, nil
}
