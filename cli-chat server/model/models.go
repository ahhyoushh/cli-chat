package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username     string             `json:"username"`
	Conversation []string           `json:"conversation"`
	Password     string             `json:"password"`
}

type Conversations struct {
	Conversation_ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	Sender   string    `json:"sender"`
	Receiver string    `json:"receiver"`
	Message  string    `json:"message"`
	Read     bool      `json:"read"`
	Time     time.Time `json:"time"`
}
