package server

import (
	"context"
	"fmt"
	"swd_project/src/db/postgresdb"
	"swd_project/src/model"
	"swd_project/src/pbs/centerpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CenterServer struct {
	centerpb.UnimplementedCenterServiceServer
}

func (*CenterServer) FindAll(ctx context.Context, req *centerpb.FindAllRequest) (*centerpb.FindAllResponse, error) {
	var centers []model.Center
	if err := postgresdb.DB.Find(&centers).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while fetching data from database : %v", err),
		)
	}
	var centertsRes = []*centerpb.Center{}
	for _, v := range centers {
		centertsRes = append(centertsRes, &centerpb.Center{
			Id:        int32(v.Model.ID),
			Website:   v.Website,
			Telephone: v.Telephone,
			Sms:       v.Sms,
		})
	}
	return &centerpb.FindAllResponse{
		Center: centertsRes,
	}, nil
}

func (*CenterServer) Find(ctx context.Context, req *centerpb.FindRequest) (*centerpb.FindResponse, error) {
	var center model.Center
	if err := postgresdb.DB.Where("id = ?", req.GetCenterId()).Find(&center).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while fetching data from database : %v", err),
		)
	}
	if center.ID == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("مرکز مشاوره یافت نشد"),
		)
	}
	return &centerpb.FindResponse{
		Center: &centerpb.Center{
			Id:        int32(center.Model.ID),
			Website:   center.Website,
			Telephone: center.Telephone,
			Sms:       center.Sms,
		},
	}, nil
}
