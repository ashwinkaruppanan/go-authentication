package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `json:"name" `
	Email    string             `json:"email" validate:"email,required" `
	Password string             `json:"password" validate:"required"`
}
