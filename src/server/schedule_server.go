package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"swd_project/src/db/postgresdb"
	"swd_project/src/model"
	"swd_project/src/pbs/schedulepb"
	"time"
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
	schedule := model.Schedule{
		Each:         int(req.GetEach()),
		ConsultantID: int(req.GetConsultantId()),
		Start:        req.GetPeriod().Start.AsTime(),
		End:          req.GetPeriod().End.AsTime(),
	}
	if err := postgresdb.DB.Create(&schedule).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("خطا هنگام ذخیره سازی وقت ها"),
		)
	}
	return &schedulepb.CreateResponse{}, nil
}

func (*ScheduleServer) FindAllSchedules(ctx context.Context, req *schedulepb.FindAllSchedulesRequest) (*schedulepb.FindAllSchedulesResponse, error) {
	var schedules = []model.Schedule{}
	if err := postgresdb.DB.Where("consultant_id = ? AND start > ?", req.GetConsultantId(), req.GetStart().AsTime()).Find(&schedules).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while fetching data from database : %v", err),
		)
	}
	fmt.Println(schedules)
	var resSchedules = []*schedulepb.Schedule{}
	for _, schedule := range schedules {
		slotsCount := int(int(schedule.End.Sub(schedule.Start).Minutes()) / schedule.Each)
		//fmt.Println(int(schedule.End.Sub(schedule.Start).Minutes()), slotsCount)
		for i := 1; i <= slotsCount; i++ {
			if schedule.Start.Add(time.Duration(i*schedule.Each)*time.Minute).Unix() > req.GetStart().AsTime().Add(24*time.Hour).Unix() {
				continue
			}
			end := timestamppb.New(schedule.Start.Add(time.Duration(i*schedule.Each) * time.Minute))
			start := timestamppb.New(schedule.Start.Add(time.Duration(i*schedule.Each) * time.Minute).Add(time.Duration(-schedule.Each) * time.Minute))
			newSchedule := &schedulepb.Schedule{
				ConsultantId: int32(schedule.ConsultantID),
				Period: &schedulepb.Period{
					Start: start,
					End:   end,
				},
				UserId: 0,
			}
			// find user for this schedule
			var reserve = model.Reserve{}
			if err := postgresdb.DB.Where("consultant_id = ? AND start = ?", schedule.ConsultantID, start.AsTime()).Find(&reserve).Error; err != nil {
				return nil, status.Errorf(
					codes.Internal,
					fmt.Sprintf("error while fetching data from database : %v", err),
				)
			}
			if reserve.UserID != 0 {
				newSchedule.UserId = int32(reserve.UserID)
			}
			resSchedules = append(resSchedules, newSchedule)
		}
	}
	return &schedulepb.FindAllSchedulesResponse{
		Schedules: resSchedules,
	}, nil
}
