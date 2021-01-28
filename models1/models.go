package models1

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

)

type HTTPUsers struct {
	Data User `json:"data"`
}


type User struct {
	ID        primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string               `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string              `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email     string              `json:"email,omitempty" bson:"email,omitempty"`
	Created   time.Time          `json:"created,omitempty" bson:"created,omitempty"`
	//HTTPUsers *HTTPUsers            `json:"httpusers" bson:"httpusers,omitempty"`
}





