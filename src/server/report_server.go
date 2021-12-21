package server

import (
	"context"
	"esme_team/src/db/mongodb"
	"esme_team/src/db/postgresdb"
	"esme_team/src/model"
	"esme_team/src/pbs/reportpb"
	"fmt"
	"time"

	"github.com/appleboy/go-fcm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ReportServer struct {
	reportpb.UnimplementedReportServiceServer
	Collection *mongo.Collection
}

func (r *ReportServer) CreateReport(ctx context.Context, req *reportpb.SensorReportRequest) (*reportpb.SensorReportResponse, error) {
	report := req.GetReport()
	var s int
	if report.GetStatus() == reportpb.Report_Default {
		s = int(reportpb.Report_LOW)
	}
	fault := model.FaultModel{
		Model: gorm.Model{
			ID: uint(report.GetFaultId()),
		},
	}
	if err := postgresdb.DB.Model(&fault).Preload("Commands").Find(&fault).Error; err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while fetch fault data : %v", err),
		)
	}
	if fault.Title == "" {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("خطا یافت نشد"),
		)
	}

	var commands []model.CommandReportModel = []model.CommandReportModel{}
	var resCommands []*reportpb.Command = []*reportpb.Command{}
	for _, v := range fault.Commands {
		commands = append(commands, model.CommandReportModel{
			CommandID: int(v.ID),
			Title:     v.Title,
			Auto:      v.Auto,
			Done:      false,
		})
		resCommands = append(resCommands, &reportpb.Command{
			CommandId: int32(v.ID),
			Title:     v.Title,
			Auto:      v.Auto,
			Done:      false,
		})

	}

	reportData := model.ReportModel{
		FaultID:   int(report.GetFaultId()),
		SensorID:  int(report.GetSensorId()),
		Status:    s,
		Tags:      report.GetTags(),
		CreatedAt: primitive.DateTime(time.Now().Minute()),
		UpdatedAt: primitive.DateTime(time.Now().Minute()),
		UserID:    int(report.GetUserId()),
		Commands:  commands,
	}

	res, err := r.Collection.InsertOne(context.Background(), reportData)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while store report data into database: %v", err),
		)
	}

	rid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error while casting stored data: %v", err),
		)
	}

	// add log
	logServer := NewLogServer()
	logData := model.LogModel{
		Operation: CreateOperation,
		ReportID:  rid.Hex(),
		SensorID:  int(report.GetSensorId()),
		UserID:    int(report.GetUserId()),
		CreatedAt: primitive.DateTime(time.Now().Minute()),
	}
	logServer.Add(logData)

	var tokens []model.TokenModel
	if err := postgresdb.DB.Find(&tokens).Error; err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("error while fetch data! : %v", err),
		)
	}
	var proiarity string
	if s == 0 || s == 1 {
		proiarity = "پایین"
	} else if s == 2 {
		proiarity = "متوسط"
	} else if s == 3 {
		proiarity = "بالا"
	}
	for _, v := range tokens {
		msg := fcm.Message{
			To: v.Token,
			Data: map[string]interface{}{
				"report_id": rid,
			},
			Notification: &fcm.Notification{
				Title: fault.Title,
				Body:  fmt.Sprintf("اهمیت: %v\nآدرس: %v", proiarity, "ساختمان صالی، تابلو برق طبقه اول"),
				Color: "#df4759",
			},
		}
		SendNotif(msg)
	}

	return &reportpb.SensorReportResponse{
		Report: &reportpb.Report{
			Id:        rid.Hex(),
			FaultId:   report.GetFaultId(),
			SensorId:  report.GetSensorId(),
			Status:    report.GetStatus(),
			Tags:      report.GetTags(),
			CreatedAt: reportData.CreatedAt.Time().String(),
			UpdatedAt: reportData.UpdatedAt.Time().String(),
			UserId:    int32(reportData.UserID),
			Commands:  resCommands,
		},
	}, nil
}

func (r *ReportServer) GetReport(ctx context.Context, req *reportpb.GetReportRequest) (*reportpb.GetReportResponse, error) {
	var report model.ReportModel
	oid, err := primitive.ObjectIDFromHex(req.GetReportId())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error while get OID: %v", err),
		)
	}
	errFind := r.Collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&report)
	if errFind == mongo.ErrNoDocuments {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("report not found!"),
		)
	}
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error while fetch data: %v", err),
		)
	}
	var commands []*reportpb.Command = []*reportpb.Command{}
	for _, v := range report.Commands {
		commands = append(commands, &reportpb.Command{
			CommandId: int32(v.CommandID),
			Title:     v.Title,
			Auto:      v.Auto,
			Done:      v.Done,
		})
	}
	res := &reportpb.GetReportResponse{
		Report: &reportpb.Report{
			Id:        req.GetReportId(),
			FaultId:   int32(report.FaultID),
			SensorId:  int32(report.SensorID),
			Status:    reportpb.Report_Status(report.Status),
			Tags:      report.Tags,
			CreatedAt: report.CreatedAt.Time().String(),
			UpdatedAt: report.UpdatedAt.Time().String(),
			UserId:    int32(report.UserID),
			Commands:  commands,
		},
	}
	return res, nil
}

func (r *ReportServer) GetUnCompletedReports(ctx context.Context, req *reportpb.GetUnCompletedReportsRequest) (*reportpb.GetUnCompletedReportsResponse, error) {
	cur, err := r.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error while fetch data: %v", err),
		)
	}
	defer cur.Close(context.Background())
	var reports []*reportpb.Report = []*reportpb.Report{}
	for cur.Next(context.Background()) {
		var result model.ReportModel
		err := cur.Decode(&result)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while fetch data: %v", err),
			)
		}
		hasUnCompletedCommand := false
		for _, v := range result.Commands {
			if !v.Done {
				hasUnCompletedCommand = true
				break
			}
		}
		if hasUnCompletedCommand || len(result.Commands) == 0 {
			var commands []*reportpb.Command = []*reportpb.Command{}
			for _, v := range result.Commands {
				commands = append(commands, &reportpb.Command{
					CommandId: int32(v.CommandID),
					Title:     v.Title,
					Auto:      v.Auto,
					Done:      v.Done,
				})
			}
			reports = append(reports, &reportpb.Report{
				Id:       result.ID.Hex(),
				FaultId:  int32(result.FaultID),
				SensorId: int32(result.SensorID),
				Status:   reportpb.Report_Status(result.Status),
				Tags:     result.Tags,
				UserId:   int32(result.UserID),
				Commands: commands,
			})
		}
	}
	return &reportpb.GetUnCompletedReportsResponse{
		Report: reports,
	}, nil
}

func (r *ReportServer) GetReportLogs(ctx context.Context, req *reportpb.GetReportLogRequest) (*reportpb.GetReportLogResponse, error) {
	oid, err := primitive.ObjectIDFromHex(req.GetReportId())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error while get OID: %v", err),
		)
	}

	cur, err := mongodb.MongoClient.Database("server_db").Collection("logs").Find(context.Background(), bson.M{})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error while fetch data: %v", err),
		)
	}
	defer cur.Close(context.Background())
	var logs []*reportpb.ReportLog = []*reportpb.ReportLog{}
	for cur.Next(context.Background()) {
		var result model.LogModel
		err := cur.Decode(&result)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while fetch data: %v", err),
			)
		}
		if result.ReportID == oid.Hex() {
			logs = append(logs, &reportpb.ReportLog{
				Operation:   result.Operation,
				UserId:      int32(result.UserID),
				ReportId:    result.ReportID,
				SendsorId:   int32(result.SensorID),
				Description: result.Description,
				CreatedAt:   result.CreatedAt.Time().String(),
			})
		}
	}
	return &reportpb.GetReportLogResponse{
		ReportLog: logs,
	}, nil
}

func SendNotif(msg fcm.Message) error {
	client, err := fcm.NewClient("")
	if err != nil {
		return err
	}
	_, errSend := client.Send(&msg)
	if errSend != nil {
		return err
	}
	return nil
}

func (r *ReportServer) GetAllLogs(ctx context.Context, req *reportpb.GetAllLogsRequest) (*reportpb.GetAllLogsResponse, error) {
	cur, err := mongodb.MongoClient.Database("server_db").Collection("logs").Find(context.Background(), bson.M{})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error while fetch data: %v", err),
		)
	}
	defer cur.Close(context.Background())
	var logs []*reportpb.ReportLog = []*reportpb.ReportLog{}
	for cur.Next(context.Background()) {
		var result model.LogModel
		err := cur.Decode(&result)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while fetch data: %v", err),
			)
		}

		logs = append(logs, &reportpb.ReportLog{
			Operation:   result.Operation,
			UserId:      int32(result.UserID),
			ReportId:    result.ReportID,
			SendsorId:   int32(result.SensorID),
			Description: result.Description,
			CreatedAt:   result.CreatedAt.Time().String(),
		})

	}
	return &reportpb.GetAllLogsResponse{
		ReportLog: logs,
	}, nil
}

func (r *ReportServer) DoReportCommand(ctx context.Context, req *reportpb.DoReportCommandRequest) (*reportpb.DoReportCommandResponse, error) {
	oid, err := primitive.ObjectIDFromHex(req.GetReportId())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error while get OID: %v", err),
		)
	}
	var report model.ReportModel
	if err := r.Collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&report); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error while fetch data: %v", err),
		)
	}
	if report.FaultID == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("گزارش یافت نشد"),
		)
	}
	commandExists := false
	var ci int
	for k, v := range report.Commands {
		if v.CommandID == int(req.GetCommandId()) {
			commandExists = true
			ci = k
			break
		}
	}
	if !commandExists {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("دستور برای گزارش یافت نشد"),
		)
	}
	var newCommands []model.CommandReportModel = []model.CommandReportModel{}
	for k, v := range report.Commands {
		if k == ci {
			newCommands = append(newCommands, model.CommandReportModel{
				CommandID: v.CommandID,
				Title:     v.Title,
				Auto:      v.Auto,
				Done:      true,
			})
		} else {
			newCommands = append(newCommands, model.CommandReportModel{
				CommandID: v.CommandID,
				Title:     v.Title,
				Auto:      v.Auto,
				Done:      v.Done,
			})
		}
	}

	r.Collection.UpdateOne(context.Background(), bson.M{"_id": oid}, bson.D{{"$set", bson.D{{"commands", newCommands}}}})

	logServer := NewLogServer()
	logData := model.LogModel{
		Operation: UpdateOperation,
		ReportID:  oid.Hex(),
		SensorID:  report.SensorID,
		UserID:    report.UserID,
		CreatedAt: primitive.DateTime(time.Now().Minute()),
	}
	logServer.Add(logData)

	return &reportpb.DoReportCommandResponse{
		Result: "دستور با موفقیت اجرا شد",
	}, nil
}

func (r *ReportServer) GetSensorReports(ctx context.Context, req *reportpb.GetSensorReportsRequest) (*reportpb.GetSensorReportsResponse, error) {
	cur, err := r.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error while fetch data: %v", err),
		)
	}
	defer cur.Close(context.Background())
	var reports []*reportpb.Report = []*reportpb.Report{}
	for cur.Next(context.Background()) {
		var result model.ReportModel
		err := cur.Decode(&result)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while fetch data: %v", err),
			)
		}
		hasUnCompletedCommand := false
		for _, v := range result.Commands {
			if !v.Done {
				hasUnCompletedCommand = true
				break
			}
		}
		if (hasUnCompletedCommand || len(result.Commands) == 0) && result.SensorID == int(req.GetSensorId()) {
			var commands []*reportpb.Command = []*reportpb.Command{}
			for _, v := range result.Commands {
				commands = append(commands, &reportpb.Command{
					CommandId: int32(v.CommandID),
					Title:     v.Title,
					Auto:      v.Auto,
					Done:      v.Done,
				})
			}
			reports = append(reports, &reportpb.Report{
				Id:       result.ID.Hex(),
				FaultId:  int32(result.FaultID),
				SensorId: int32(result.SensorID),
				Status:   reportpb.Report_Status(result.Status),
				Tags:     result.Tags,
				UserId:   int32(result.UserID),
				Commands: commands,
			})
		}
	}
	return &reportpb.GetSensorReportsResponse{
		Reports: reports,
	}, nil
}

func (r *ReportServer) GetSensorsReportsCount(ctx context.Context, req *reportpb.GetSensorsReportsCountRequest) (*reportpb.GetSensorsReportsCountResponse, error) {
	// id -> counter
	counter := map[int]int{}
	cur, err := r.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Error while fetch data: %v", err),
		)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var result model.ReportModel
		err := cur.Decode(&result)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while fetch data: %v", err),
			)
		}
		counter[int(result.SensorID)]++
	}
	for i := 1; i <= 5; i++ {
		if _ , exists := counter[i]; exists == false {
			counter[i] = 0
		}
	}
	var resArr []*reportpb.GetSensorsReportsCountResponse_SensorIdWithReportCount = []*reportpb.GetSensorsReportsCountResponse_SensorIdWithReportCount{}
	for k, v := range counter {
		resArr = append(resArr, &reportpb.GetSensorsReportsCountResponse_SensorIdWithReportCount{
			SensorId: int32(k),
			Count:    int32(v),
		})
	}
	
	return &reportpb.GetSensorsReportsCountResponse{
		SensorIdWithReportCount: resArr,
	}, nil
}
