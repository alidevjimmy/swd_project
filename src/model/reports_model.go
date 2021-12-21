package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ReportModel struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	FaultID   int                  `bson:"fault_id"`
	SensorID  int                  `bson:"sensor_id"`
	Status    int                  `bson:"status"`
	Tags      []string             `bson:"tags"`
	Commands  []CommandReportModel `bson:"commands"`
	UserID    int                  `bson:"user_id"`
	CreatedAt primitive.DateTime   `bson:"created_at"`
	UpdatedAt primitive.DateTime   `bson:"updated_at"`
}

type CommandReportModel struct {
	CommandID int    `bson:"command_id"`
	Title     string `bson:"title"`
	Auto      bool   `bson:"auto"`
	Done      bool   `bson:"done"`
}
