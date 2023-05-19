package main

import (
	"aapi/config"
	"aapi/internal/interceptors"
	"aapi/pkg/logger"
	appsvc "aapi/protos/v1"
	logging_service "aapi/services/logging/protos/log/v1"
	user_sv "aapi/services/user/protos/user/v1"
	"aapi/shared/constants"
	"context"
	"fmt"
	"net"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v8"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/minio/minio-go/v7"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type appServer struct {
	logger             logger.Logger
	cfg                *config.Config
	tracer             opentracing.Tracer
	redisClient        *redis.Client
	redisReplicaClient *redis.Client
	userApiClient      user_sv.UserClient
	loggingApiClient   logging_service.LoggingClient
	kafkaProducer      *kafka.Producer
	minioClient        *minio.Client
	appsvc.UnimplementedUserServer
	appsvc.UnimplementedLoggingServer
}

type AppServer interface {
	//User service
	Enroll(ctx context.Context, model *appsvc.EnrollRequest) (*appsvc.EnrollResponse, error)
	Verify(ctx context.Context, model *appsvc.VerifyRequest) (*appsvc.VerifyResponse, error)
	GetUser(ctx context.Context, model *appsvc.GetUserRequest) (*appsvc.GetUserResponse, error)
	CountUser(ctx context.Context, model *appsvc.CountUserRequest) (*appsvc.CountUserResponse, error)
	SearchUser(ctx context.Context, model *appsvc.SearchUserRequest) (*appsvc.SearchUserResponse, error)
	Update(ctx context.Context, model *appsvc.UpdateRequest) (*appsvc.UpdateResponse, error)
	Delete(ctx context.Context, model *appsvc.DeleteRequest) (*appsvc.DeleteResponse, error)

	// Logging
	AddLog(ctx context.Context, model *appsvc.AddLogRequest) (*appsvc.AddLogResponse, error)
	GetLog(ctx context.Context, model *appsvc.QueryRequest) (*appsvc.QueryResponse, error)
	GetLogs(ctx context.Context, model *appsvc.GetLogsRequest) (*appsvc.GetLogsResponse, error)

	Run(port int) error
}

// NewService func initializes a service
func NewServer(
	logger logger.Logger,
	cfg *config.Config,
	userServiceApiCnn *grpc.ClientConn,
	loggingServiceApiCnn *grpc.ClientConn,
	trace opentracing.Tracer,
	redisClient *redis.Client,
	redisReplicaClient *redis.Client,
	kProducer *kafka.Producer,
	minioClient *minio.Client,
) AppServer {
	return &appServer{
		logger:             logger,
		cfg:                cfg,
		userApiClient:      user_sv.NewUserClient(userServiceApiCnn),
		loggingApiClient:   logging_service.NewLoggingClient(loggingServiceApiCnn),
		tracer:             trace,
		redisClient:        redisClient,
		redisReplicaClient: redisReplicaClient,
		kafkaProducer:      kProducer,
		minioClient:        minioClient,
	}
}

// Run starts the server
func (s *appServer) Run(port int) error {

	auditLogInterceptors := interceptors.NewAuditLogInterceptor(s.logger, s.cfg, s.kafkaProducer)

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(auditLogInterceptors.Unary())),
		grpc.MaxRecvMsgSize(constants.MAX_MESSAGE_LENGTH),
		grpc.MaxSendMsgSize(constants.MAX_MESSAGE_LENGTH),
	)

	appsvc.RegisterUserServer(srv, s)
	appsvc.RegisterLoggingServer(srv, s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.logger.Infof("failed to listen: %v", err)
	}
	return srv.Serve(lis)
}
