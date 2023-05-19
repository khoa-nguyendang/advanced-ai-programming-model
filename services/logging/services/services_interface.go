package services

import (
	"aapi/config"
	"aapi/pkg/logger"
	logging_service "aapi/services/logging/protos/log/v1"
	"aapi/services/logging/repositories"
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-co-op/gocron"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	USER_SERVICE_VERIFICATION_TOPIC = "user-service-verification-topic"
)

type service struct {
	logger           logger.Logger
	tracer           opentracing.Tracer
	cfg              *config.Config
	kafkaConsumer    *kafka.Consumer
	auditLogConsumer *kafka.Consumer
	scheduler        *gocron.Scheduler
	logRepository    repositories.LogRepository
	mongoClient      *mongo.Client
	logging_service.UnimplementedLoggingServer
}

type LoggingService interface {
	AddLog(ctx context.Context, model *logging_service.AddLogRequest) (*logging_service.AddLogResponse, error)
	GetLog(ctx context.Context, model *logging_service.QueryRequest) (*logging_service.QueryResponse, error)
	GetLogs(ctx context.Context, model *logging_service.GetLogsRequest) (*logging_service.GetLogsResponse, error)
	ScheduleReport(ctx context.Context, model *logging_service.ScheduleReportRequest) (*logging_service.ScheduleReportResponse, error)
	Run(port int) error
}

// NewService func initializes a service
func NewService(logger logger.Logger,
	log_repository repositories.LogRepository,
	trace opentracing.Tracer,
	cfg *config.Config,
	kafka_consumer *kafka.Consumer,
	auditLogConsumer *kafka.Consumer,
	scheduler *gocron.Scheduler,
	mongoClient *mongo.Client,
) LoggingService {
	return &service{
		logger:           logger,
		logRepository:    log_repository,
		tracer:           trace,
		cfg:              cfg,
		kafkaConsumer:    kafka_consumer,
		auditLogConsumer: auditLogConsumer,
		scheduler:        scheduler,
		mongoClient:      mongoClient,
	}
}
