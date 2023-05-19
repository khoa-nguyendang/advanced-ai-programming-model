package repositories

import (
	"aapi/pkg/logger"
	lsv "aapi/services/logging/protos/log/v1"
	"aapi/shared/constants"
	"aapi/shared/entities"
	viewmodels "aapi/shared/view_models"
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
)

type logRepo struct {
	db          *sqlx.DB
	cache       *redis.Client
	logger      logger.Logger
	mongoClient *mongo.Client
}

type LogRepository interface {
	GetLog(ctx context.Context, start_date int64, end_date int64, device_uuids []string, user_ids []string, company_code string) []*entities.Log
	GetLogs(ctx context.Context, model *lsv.GetLogsRequest) (*lsv.GetLogsResponse, error)
	SaveLog(app_user_id int64, activity constants.Activity, device_uuid string, full_message string, date int64) error
	SaveAuditLog(auditLogJson string, mdc viewmodels.MongoDbConfig) error
}

// NewRepository func initializes a service
func NewRepository(db *sqlx.DB,
	logger logger.Logger,
	cache *redis.Client,
	mongoClient *mongo.Client) LogRepository {
	return &logRepo{db: db, logger: logger, cache: cache, mongoClient: mongoClient}
}
