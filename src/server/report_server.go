package server

import (
	"context"
	"swd_project/src/pbs/reportpb"
)

type ReportServer struct {
	reportpb.UnimplementedReportServiceServer
}

func (*ReportServer) CreateReport(ctx context.Context,req *reportpb.CreateReportRequest) (*reportpb.CreateReportResponse, error) {
	return nil, nil
}
func (*ReportServer) UserOpenReports(ctx context.Context,req *reportpb.UserOpenReportsRequest) (*reportpb.UserOpenReportsResponse, error) {
	// only reports
	return nil, nil
}