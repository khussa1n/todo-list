package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tasks struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"title"`
	ActiveAt string             `json:"activeAt" bson:"activeAt"`
	Status   string             `json:"status" bson:"status"`
}
