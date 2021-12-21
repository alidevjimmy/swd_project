package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type LogModel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Operation   string             `bson:"operation"`
	ReportID    string             `bson:"report_id"`
	SensorID    int                `bson:"sensor_id"`
	UserID      int                `bson:"user_id"`
	Description string             `bson:"description"`
	CreatedAt   primitive.DateTime `bson:"created_at"`
}
