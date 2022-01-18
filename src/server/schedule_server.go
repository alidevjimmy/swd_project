package server

import (
	"context"
	"encoding/json"
	"fmt"
	"swd_project/src/db/postgresdb"
	"swd_project/src/model"
	"swd_project/src/pbs/schedulepb"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ScheduleServer struct {
	schedulepb.UnimplementedScheduleServiceServer
}

func (*ScheduleServer) Create(ctx context.Context, req *schedulepb.CreateRequest) (*schedulepb.CreateResponse, error) {
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
	var periods = model.Periods{}
	for _, period := range req.GetPeriods() {
		periods = append(periods, struct {
			Start time.Time `json:"start"`
			End   time.Time `json:"end"`
		}{
			Start: period.Start.AsTime(),
			End:   period.End.AsTime(),
		})
	}
	schedule := model.Schedule{
		Each:         int(req.GetEach()),
		ConsultantID: int(req.GetConsultantId()),
	}

	jsonByte, _ := json.Marshal(periods)
	schedule.Periods = jsonByte
	if err := postgresdb.DB.Create(&schedule).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("خطا هنگام ذخیره سازی وقت ها"),
		)
	}
	return &schedulepb.CreateResponse{}, nil
}

func (*ScheduleServer) FindAllSchedules(ctx context.Context, req *schedulepb.FindAllSchedulesRequest) (*schedulepb.FindAllSchedulesResponse, error) {
	// var schedules = []model.Schedule{}
	// if err := postgresdb.DB.Where("date >= ?", time.Now()).Find(&schedules).Error; err != nil {
	// 	return nil, status.Errorf(
	// 		codes.NotFound,
	// 		fmt.Sprintf("خطا هنگام استخراج وقت های آزاد"),
	// 	)
	// }
	// var schedulesRes = []*schedulepb.Schedule{}
	// for _, v := range schedules {
	// 	if v.Date == time.Now() {
	// 		if time.Now().Local().Hour() >= v.FinishHour {
	// 			continue
	// 		}
	// 	}
	// 	date := timestamppb.New(v.Date)
	// 	schedulesRes = append(schedulesRes, &schedulepb.Schedule{
	// 		Id:           int32(v.ID),
	// 		Date:         date,
	// 		StartHour:    int32(v.StartHour),
	// 		FinishHour:   int32(v.FinishHour),
	// 		ConsultantId: int32(v.ConsultantID),
	// 	})
	// }
	// return &schedulepb.FindAllOpenSchedulesResponse{
	// 	Schedules: schedulesRes,
	// }, nil
	return &schedulepb.FindAllSchedulesResponse{}, nil
}
