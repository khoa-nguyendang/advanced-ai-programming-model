package interceptors

import (
	"aapi/config"
	"aapi/pkg/kafka_srv"
	"aapi/pkg/logger"
	"aapi/shared/constants"
	"aapi/shared/entities"
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/grpc"
)

type AuditLogManager struct {
	logger        logger.Logger
	cfg           *config.Config
	kafkaProducer *kafka.Producer
	ignoredAPIs   map[string]bool
}

func NewAuditLogInterceptor(logger logger.Logger, cfg *config.Config, kProducer *kafka.Producer) *AuditLogManager {
	return &AuditLogManager{logger: logger, cfg: cfg, kafkaProducer: kProducer, ignoredAPIs: ignoredAPIs()}
}

func (alm *AuditLogManager) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// won't handle API for verify
		if _, found := alm.ignoredAPIs[info.FullMethod]; found {
			return handler(ctx, req)
		}
		alm.logger.Infof("AuditLogManager.Unary.handling api: %v", info.FullMethod)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go alm.RecordRequest(&wg, req, int(constants.DATAUPDATED_COMMON), info.FullMethod, info.FullMethod, info.FullMethod)
		wg.Wait()
		return handler(ctx, req)
	}
}

func (alm *AuditLogManager) RecordRequest(wg *sync.WaitGroup, payload interface{}, updateType int, from, entity, service string) {
	defer wg.Done()
	event := &entities.AuditLog{
		Entity:          entity,
		Service:         service,
		ModifiedBy:      from,
		DataUpdatedType: updateType,
		ModifiedAt:      time.Now().UTC().UnixMilli(),
		Payload:         payload,
	}
	message, _ := json.Marshal(event)
	alm.logger.Infof("AuditLogManager.RecordRequest: %v", string(message))
	kafka_srv.SendMessage(alm.kafkaProducer, constants.KAFKATOPIC_DATA_UPDATED, string(message))
}

// ignoredAPIs don't need to listen thoses GET request from client
func ignoredAPIs() map[string]bool {
	result := make(map[string]bool)

	result["/appsvc.Authentication/GetToken"] = true
	result["/appsvc.Authentication/AdministratorGetToken"] = true
	result["/appsvc.Administrator/GetAdministratorInfo"] = true
	result["/appsvc.Administrator/ContactAdministrator"] = true
	result["/appsvc.Administrator/GetAdministratorList"] = true
	result["/appsvc.Company/GetCompanies"] = true
	result["/appsvc.Company/GetGroups"] = true
	result["/appsvc.Logging/AddLog"] = true
	result["/appsvc.Logging/GetLog"] = true
	result["/appsvc.Logging/GetLogs"] = true
	result["/appsvc.Location/GetLocation"] = true
	result["/appsvc.Device/GetDevices"] = true
	result["/appsvc.Device/GetDevice"] = true
	result["/appsvc.Device/GetMessageByUuid"] = true
	result["/appsvc.User/Enroll"] = true
	result["/appsvc.User/Verify"] = true
	result["/appsvc.User/GetUser"] = true
	result["/appsvc.User/SearchUser"] = true
	result["/appsvc.User/GetGroups"] = true
	result["/appsvc.User/GetRoles"] = true
	result["/appsvc.User/Verify"] = true
	result["/appsvc.User/CountUser"] = true
	return result
}
