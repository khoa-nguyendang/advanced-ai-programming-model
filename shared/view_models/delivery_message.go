package viewmodels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeliveryMessage struct {
	Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MessageUuid string             `json:"message_uuid,omitempty" bson:"message_uuid,omitempty"`
	MessageType string             `json:"message_type,omitempty" bson:"message_type,omitempty"`
}
