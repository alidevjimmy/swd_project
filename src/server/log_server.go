package server

import (
	"context"
	"esme_team/src/db/mongodb"
	"esme_team/src/model"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CreateOperation = "CREATE"
	UpdateOperation = "UPDATE"
	DeleteOperation = "DELETE"
)

type LogServer struct{}

func NewLogServer() LogServer {
	return LogServer{}
}

func (*LogServer) Add(log model.LogModel) error {
	_, err := mongodb.MongoClient.Database("server_db").Collection("logs").InsertOne(context.Background(), log)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while inserting log into database : %v", err),
		)
	}
	return nil
}
