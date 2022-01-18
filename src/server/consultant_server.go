package server

import (
	"context"
	"fmt"
	"swd_project/src/db/postgresdb"
	"swd_project/src/model"
	"swd_project/src/pbs/consultantpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConsultantServer struct {
	consultantpb.UnimplementedConsultantServiceServer
}

func (*ConsultantServer) FindAllConsultants(ctx context.Context, req *consultantpb.FindAllConsultantsRequest) (*consultantpb.FindAllConsultantsResponse, error) {
	var consultants []model.Consultant
	if err := postgresdb.DB.Find(&consultants).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while fetching data from database : %v", err),
		)
	}
	var consultanttsRes = []*consultantpb.Consultant{}
	for _, v := range consultants {
		consultanttsRes = append(consultanttsRes, &consultantpb.Consultant{
			Id:          int32(v.Model.ID),
			Name:        v.Name,
			Family:      v.Family,
			Description: v.Description,
			Phone:       v.Phone,
		})
	}
	return &consultantpb.FindAllConsultantsResponse{
		Consultants: consultanttsRes,
	}, nil
}

func (*ConsultantServer) FindConsultant(ctx context.Context, req *consultantpb.FindConsultantRequest) (*consultantpb.FindConsultantResponse, error) {
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
	return &consultantpb.FindConsultantResponse{
		Consultant: &consultantpb.Consultant{
			Id:          int32(consultant.Model.ID),
			Name:        consultant.Name,
			Family:      consultant.Family,
			Description: consultant.Description,
			Phone:       consultant.Phone,
		},
	}, nil
}
