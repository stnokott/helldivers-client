package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PrimitiveTime(t time.Time) primitive.Timestamp {
	return primitive.Timestamp{T: uint32(t.Unix())}
}
