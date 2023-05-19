package kafka_srv

import (
	"aapi/config"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewKafkaProducer(c *config.Config) (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": c.Kafka.Server})
	return producer, err
}

func NewKafkaConsumer(c *config.Config, group string) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  c.Kafka.Server,
		"group.id":           group,
		"auto.offset.reset":  kafka.OffsetBeginning.String(),
		"enable.auto.commit": false,
	})
	return consumer, err
}

func NewKafkaConsumerAuditLog(c *config.Config, group string) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  c.Kafka.AuditLogServer,
		"group.id":           group,
		"auto.offset.reset":  kafka.OffsetBeginning.String(),
		"enable.auto.commit": false,
	})
	return consumer, err
}

func SendMessage(producer *kafka.Producer, topic string, message string) {
	// Produce messages to topic (asynchronously)
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)

	if err != nil {
		log.Default().Printf("SendMessage.err: %v", err)
	} else {
		log.Default().Printf("Sent message: %v --- topic: %v", message, topic)
	}

}
