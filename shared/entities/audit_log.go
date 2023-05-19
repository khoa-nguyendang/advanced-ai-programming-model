package entities

type AuditLog struct {
	ID              string      `json:"id" bson:"_id"`
	Entity          string      `json:"entity" bson:"entity"`
	Service         string      `json:"service" bson:"service"`
	ModifiedBy      string      `json:"modified_by" bson:"modified_by"`
	ModifiedAt      int64       `json:"modified_at" bson:"modified_at"`
	DataUpdatedType int         `json:"data_updated_type" bson:"data_updated_type"`
	Payload         interface{} `json:"payload" bson:"payload"`
}
