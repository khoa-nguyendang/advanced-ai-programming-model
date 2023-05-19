package viewmodels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommonMessage struct {
	Id              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MessageUuid     string             `json:"message_uuid,omitempty" bson:"message_uuid,omitempty"`
	MessageType     string             `json:"message_type,omitempty" bson:"message_type,omitempty"`
	From            string             `json:"from,omitempty" bson:"from,omitempty"`
	To              string             `json:"to,omitempty" bson:"to,omitempty"`
	RecipientAlias  string             `json:"recipient_alias,omitempty" bson:"recipient_alias,omitempty"`
	SenderAlias     string             `json:"sender_alias,omitempty" bson:"sender_alias,omitempty"`
	CompanyCode     string             `json:"company_code,omitempty" bson:"company_code,omitempty"`
	Topic           string             `json:"topic,omitempty"  bson:"topic,omitempty"`
	Body            string             `json:"body,omitempty"  bson:"body,omitempty"`
	CreatedDate     int64              `json:"created_date,omitempty" bson:"created_date,omitempty"`
	LastModified    int64              `json:"last_modified,omitempty" bson:"last_modified,omitempty"`
	Acknowledge     bool               `json:"acknowledge,omitempty"  bson:"acknowledge,omitempty"`
	AcknowledgeDate int64              `json:"acknowledge_date,omitempty"  bson:"acknowledge_date,omitempty"`
}
