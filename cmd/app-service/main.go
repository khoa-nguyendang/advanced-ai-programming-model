package main

import (
	"aapi/config"
	"aapi/pkg/kafka_srv"
	"aapi/pkg/logger"
	"aapi/pkg/minio_srv"
	"aapi/pkg/mongodb"
	"aapi/pkg/mysql"
	"aapi/pkg/trace"
	"aapi/shared/constants"
	viewmodels "aapi/shared/view_models"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	logging_repo "aapi/services/logging/repositories"
	logging_sv "aapi/services/logging/services"

	user_repo "aapi/services/user/repositories"
	user_sv "aapi/services/user/services"

	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
}

type server interface {
	Run(port int) error
}

// Server

func main() {
	var (
		port           = flag.Int("port", 8080, "The service port")
		jaegeraddr     = flag.String("jaeger", "jaeger:6831", "Jaeger address")
		configfilepath = flag.String("config-path", "", "The absolute path to configuration file")
	)
	flag.Parse()
	log.Default().Printf("file setting path: %v --- default path: %v \n", *configfilepath, os.Getenv("config"))
	// load configuration
	var configPath string
	if configfilepath != nil && *configfilepath != "" {
		configPath = *configfilepath
	} else {
		//load default
		configPath = config.GetConfigPath(os.Getenv("config"))
		log.Printf("Loadding default config path : %v \n", configPath)
	}
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	// initial logger
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v",
		cfg.Server.AppVersion,
		cfg.Logger.Level,
		cfg.Server.Mode,
		cfg.Server.SSL,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.Server.AppVersion)
	appLogger.Infof("mysql config: %#v", cfg.Mysql)
	appLogger.Infof("faiss config: %#v", cfg.AppFaiss)
	appLogger.Infof("engine config: %#v", cfg.AppEngine)

	// Jaeger trace
	tracer, err := trace.New("appTracer", *jaegeraddr)
	if err != nil {
		appLogger.Infof("trace new error: %v", err)
	}

	var srv server
	var cmd = os.Args[1]
	new_port := os.Args[2]

	//Redis
	redisMasterDb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//Redis replica
	redisReplicas := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Mysql
	mysqlDb, err := mysql.NewMysqlDB(cfg)
	if err != nil {
		appLogger.Infof("MySql init: %s", err)
	} else {
		appLogger.Infof("MySql connected: %#v", mysqlDb.Stats())
	}
	defer mysqlDb.Close()

	switch cmd {
	case "aapiserver":
		var err error
		kafkaProducer, err := kafka_srv.NewKafkaProducer(cfg)
		if err != nil {
			appLogger.Infof("Kafka init 1: %s", err)
		}

		kafka_srv.SendMessage(kafkaProducer, constants.KAFKATOPIC_DATA_UPDATED, "")
		kafka_srv.SendMessage(kafkaProducer, constants.KAFKATOPIC_USER_SERVICE_VERIFICATION, "")

		err = mysql.RunMigration(cfg)
		if err != nil {
			appLogger.Infof("MySql RunMigration: %s", err)
		}

		minioClient, err := minio_srv.NewMinIOClient(cfg)
		if err != nil {
			appLogger.Errorf("Unable to initial Minio Client: %v", err)
		}
		srv = NewServer(
			appLogger,
			cfg,
			dial(cfg.UserService, tracer),
			dial(cfg.LoggingService, tracer),
			tracer,
			redisMasterDb,
			redisReplicas,
			kafkaProducer,
			minioClient,
		)

	case "user-service":
		minio, err := minio_srv.NewMinIOClient(cfg)

		if err != nil {
			appLogger.Infof("MinIO init: %s", err)
		}

		kafka_producer, err := kafka_srv.NewKafkaProducer(cfg)

		if err != nil {
			appLogger.Infof("Kafka init: %s", err)
		}

		userRepo := user_repo.NewRepository(mysqlDb, appLogger, redisMasterDb)
		srv = user_sv.NewService(
			appLogger,
			cfg,
			dial(cfg.AppEngine.URL, tracer),
			dial(cfg.AppLightweightEngine.URL, tracer),
			dial(cfg.AppFaiss.URL, tracer),
			userRepo,
			tracer,
			minio,
			kafka_producer,
		)

	case "logging-service":
		kafka_consumer, err := kafka_srv.NewKafkaConsumer(cfg, "logging_service")
		if err != nil {
			appLogger.Errorf("Kafka NewKafkaConsumer: %s \n", err)
		}
		auditLogConsumer, err := kafka_srv.NewKafkaConsumerAuditLog(cfg, "audit_log")
		if err != nil {
			appLogger.Errorf("Kafka NewKafkaConsumerAuditLog: %s \n", err)
		}
		scheduler := gocron.NewScheduler(time.Local)
		mongoClient, err := mongodb.GetClient(context.Background(), getMongodbConfig(cfg))
		if err != nil {
			appLogger.Infof("Kafka init: %s \n", err)
		}

		logRepository := logging_repo.NewRepository(mysqlDb, appLogger, redisMasterDb, mongoClient)
		srv = logging_sv.NewService(appLogger, logRepository, tracer, cfg, kafka_consumer, auditLogConsumer, scheduler, mongoClient)
	default:
		appLogger.Infof("unknown cmd: %s", cmd)
	}

	i, err := strconv.Atoi(new_port)
	if err != nil {
		if err := srv.Run(*port); err != nil {
			appLogger.Infof("run %s error: %v", cmd, err)
		}
	} else {
		if err := srv.Run(i); err != nil {
			appLogger.Infof("run %s error: %v", cmd, err)
		}
	}

}

func dial(addr string, t opentracing.Tracer) *grpc.ClientConn {
	log.Printf("addr: %#v", addr)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(t)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(constants.MAX_MESSAGE_LENGTH), grpc.MaxCallSendMsgSize(constants.MAX_MESSAGE_LENGTH)),
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		panic(fmt.Sprintf("ERROR: dial error: %v", err))
	}

	return conn
}

func getMongodbConfig(cfg *config.Config) viewmodels.MongoDbConfig {
	url := "mongodb://" +
		cfg.Mongodb.MongodbUser +
		":" + cfg.Mongodb.MongodbPassword +
		"@" + cfg.Mongodb.MongodbHost +
		":" + cfg.Mongodb.MongodbPort +
		"/" + cfg.Mongodb.MongodbDbname +
		"?authSource=" + cfg.Mongodb.MongodbDbname

	mongoConfig := viewmodels.MongoDbConfig{
		UserName:   cfg.Mongodb.MongodbUser,
		Password:   cfg.Mongodb.MongodbPassword,
		Url:        url,
		Database:   cfg.Mongodb.MongodbDbname,
		Collection: "AuditLog",
	}
	return mongoConfig
}
