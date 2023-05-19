package services

import (
	"aapi/pkg/mongodb"
	"aapi/shared/constants"
	viewmodels "aapi/shared/view_models"
	"time"

	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *service) BackgroundRun() {
	sch := gocron.NewScheduler(time.UTC)

	sch.Every(15).Seconds().Do(s.QueryUnAckMessage)
	sch.Every(1).Day().At("01:00").Do(s.ClearOldMessages)

	// you can start running the scheduler in two different ways:
	// starts the scheduler asynchronously
	sch.StartAsync()
}

// QueryUnAckMessage get message is not ack, not processed, and last_modified is older than 30s
func (s *service) QueryUnAckMessage() {
	last30secs := time.Now().UTC().Add(-15 * time.Second).UnixMilli()
	filter := bson.D{{Key: "$and", Value: []interface{}{
		bson.D{{Key: "last_modified", Value: bson.M{"$lt": last30secs}}},
		bson.D{{Key: "processed", Value: bson.M{"$exists": false}}},
		bson.D{{Key: "acknowledge", Value: bson.M{"$exists": false}}},
	}}}
	result := make([]*viewmodels.CommonMessage, 0)
	err := mongodb.GetMany(s.mongoClient, s.getMongodbConfig(), constants.MGC_NOTIFICATION_MESSAGES, 20, 0, filter, result)
	if err != nil {
		s.logger.Errorf("QueryUnAckMessage.err: %v \n", err)
	}
	s.logger.Infof("QueryUnAckMessage result set count : %v \n", len(result))
}

func (s *service) PushUnAckMessageToWeb(messages []*viewmodels.CommonMessage) {

}

func (s *service) markMessageProcessed(messages []*viewmodels.CommonMessage) {

}

// ClearOldMessages clear message in  collection name "notification_messages" older than a month from current running
func (s *service) ClearOldMessages() {

}
