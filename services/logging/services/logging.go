package services

import (
	lgsv "aapi/services/logging/protos/log/v1"
	"aapi/shared/constants"
	"aapi/shared/entities"
	viewmodels "aapi/shared/view_models"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"google.golang.org/grpc"
)

var company_code string = ""
var device_id string = ""

// Run starts the server
func (s *service) Run(port int) error {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(s.tracer),
		),
	)
	lgsv.RegisterLoggingServer(srv, s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.logger.Infof("failed to listen: %v", err)
	}
	// Start consumer thread
	if s.kafkaConsumer != nil {
		go s.handleVerifyConsumer()
	}
	go s.handleAuditLogConsumer()
	return srv.Serve(lis)
}

func (s *service) ScheduleReport(ctx context.Context, model *lgsv.ScheduleReportRequest) (*lgsv.ScheduleReportResponse, error) {
	s.logger.Infof("ScheduleReport")

	var res lgsv.ScheduleReportResponse
	time_count := 1
	if model.JobTimeCount > 0 {
		time_count = int(model.JobTimeCount)
	}
	switch model.JobTime {
	case lgsv.JobTime_EVERY_HOUR:
		s.scheduler.Every(time_count).Hours().Do(s.reportLog, ctx, model.JobTime, model.ReportType, model.Recipients)
		break
	case lgsv.JobTime_EVERY_DAY:
		s.scheduler.Every(time_count).Days().Do(s.reportLog, ctx, model.ReportType, model.Recipients)
		break
	case lgsv.JobTime_EVERY_WEEK:
		s.scheduler.Every(time_count).Weeks().Do(s.reportLog, ctx, model.ReportType, model.Recipients)
		break
	case lgsv.JobTime_EVERY_MONTH:
		s.scheduler.Every(time_count).Months().Do(s.reportLog, ctx, model.ReportType, model.Recipients)
		break
	case lgsv.JobTime_EVERY_SUNDAY:
		s.scheduler.Every(time_count).Sunday().Do(s.reportLog, ctx, model.ReportType, model.Recipients)
		break
	case lgsv.JobTime_EVERY_MONDAY:
		s.scheduler.Every(time_count).Monday().Do(s.reportLog, ctx, model.ReportType, model.Recipients)
		break
	case lgsv.JobTime_EVERY_TUESDAY:
		s.scheduler.Every(time_count).Tuesday().Do(s.reportLog, ctx, model.ReportType, model.Recipients)
		break
	case lgsv.JobTime_EVERY_WEDNESDAY:
		s.scheduler.Every(time_count).Wednesday().Do(s.reportLog, ctx, model.ReportType, model.Recipients)
		break
	case lgsv.JobTime_EVERY_THURSDAY:
		s.scheduler.Every(time_count).Thursday().Do(s.reportLog, ctx, model.ReportType, model.Recipients)
		break
	case lgsv.JobTime_EVERY_FRIDAY:
		s.scheduler.Every(time_count).Friday().Do(s.reportLog, ctx, model.ReportType, model.Recipients)
		break
	case lgsv.JobTime_EVERY_SATURDAY:
		s.scheduler.Every(time_count).Saturday().Do(s.reportLog, model.ReportType, model.Recipients)
		break
	default:
		s.logger.Infof("Cannot schedule report")
	}

	s.scheduler.StartAsync()

	res.Data = "OK"
	return &res, nil
}

// AddLog add verify log for tracking purpose when device is running by offline mode
func (s *service) AddLog(ctx context.Context, model *lgsv.AddLogRequest) (*lgsv.AddLogResponse, error) {
	s.logger.Infof("AddLog: %v", model)
	if model.Body == "" {
		return &lgsv.AddLogResponse{
			Code:    http.StatusBadRequest,
			Message: "Body is empty",
			Error:   "MODEL_INVALID",
			Status:  false,
		}, nil
	}

	jsonData, err := json.Marshal(model)
	if err != nil {
		return &lgsv.AddLogResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   "MODEL_INVALID",
			Status:  false,
		}, nil
	}
	err = s.logRepository.SaveLog(0, constants.VERIFY, model.DeviceUuid, string(jsonData), time.Now().UTC().UnixMilli())
	if err != nil {
		return &lgsv.AddLogResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   "FAILED_TO_ADD",
			Status:  false,
		}, nil
	}
	return &lgsv.AddLogResponse{
		Code:    http.StatusOK,
		Message: "Add Log successfully",
		Error:   "",
		Status:  true,
	}, nil
}

func (s *service) GetLog(ctx context.Context, model *lgsv.QueryRequest) (*lgsv.QueryResponse, error) {
	s.logger.Infof("GetLog")

	records := s.logRepository.GetLog(ctx, model.StartDate, model.EndDate, model.DeviceUuids, model.UserIds, company_code)
	var logs []*lgsv.Log
	csvData := [][]string{
		{"Company", "User ID", "User Name", "DeviceUuid", "Activity", "Date"},
	}

	for _, log := range records {
		logs = append(logs, &lgsv.Log{
			Id:          log.Id,
			UserId:      log.UserId,
			CompanyCode: log.CompanyCode,
			UserName:    log.UserName,
			UserState:   lgsv.UserState(log.UserState),
			DeviceUuid:  log.DeviceUuid,
			Activity:    int64(log.Activity),
			Date:        log.Date,
		})

		if model.SaveToFile {
			csvData = append(csvData, []string{log.CompanyCode, log.UserId, log.UserName, log.DeviceUuid, strconv.Itoa(int(log.Activity)), time.Unix(log.Date, 0).Local().Format(time.RFC3339)})
		}
	}

	var byteBuf []byte
	if model.SaveToFile {
		s.logger.Infof("Save to csv file")
		byteBuf = writeToCSV(csvData)

		// Temporary code to debug
		f, _ := os.Create("/tmp/data.csv")
		defer f.Close()
		n2, _ := f.Write(byteBuf)
		s.logger.Infof("Wrote %d bytes\n", n2)
		f.Sync()
	}
	return &lgsv.QueryResponse{
		Logs:       logs,
		ReportFile: byteBuf,
	}, nil
}

func (s *service) GetLogs(ctx context.Context, model *lgsv.GetLogsRequest) (*lgsv.GetLogsResponse, error) {
	s.logger.Infof("GetLogs")

	company_code := ""
	model.CompanyCode = company_code
	var err error
	data, err := s.logRepository.GetLogs(ctx, model)

	if err != nil {
		s.logger.Errorf("Cannot get company_code")
		return &lgsv.GetLogsResponse{
			Code:       http.StatusBadRequest,
			Message:    err.Error(),
			Error:      "FAILED_TO_GET_LOGS",
			Data:       []*lgsv.Log{},
			ReportFile: []byte{},
			Total:      0,
		}, err
	}

	var logs []*lgsv.Log = make([]*lgsv.Log, 0)
	csvData := [][]string{
		{"Company", "User ID", "User Name", "DeviceUuid", "Activity", "Date"},
	}

	for _, log := range data.Data {
		logs = append(logs, &lgsv.Log{
			Id:          log.Id,
			UserId:      log.UserId,
			CompanyCode: log.CompanyCode,
			UserName:    log.UserName,
			UserState:   lgsv.UserState(log.UserState),
			DeviceUuid:  log.DeviceUuid,
			Activity:    int64(log.Activity),
			Date:        log.Date,
		})

		if model.SaveToFile {
			csvData = append(csvData, []string{log.CompanyCode, log.UserId, log.UserName, log.DeviceUuid, strconv.Itoa(int(log.Activity)), time.Unix(log.Date, 0).Local().Format(time.RFC3339)})
		}
	}

	var byteBuf []byte
	if model.SaveToFile {
		s.logger.Infof("Save to csv file")
		byteBuf = writeToCSV(csvData)

		// Temporary code to debug
		f, _ := os.Create("/tmp/data.csv")
		defer f.Close()
		n2, _ := f.Write(byteBuf)
		s.logger.Infof("Wrote %d bytes\n", n2)
		f.Sync()
	}
	return &lgsv.GetLogsResponse{
		Data:       logs,
		ReportFile: byteBuf,
		Total:      data.Total,
		Code:       http.StatusOK,
		Message:    "",
		Error:      "",
	}, nil
}

func (s *service) reportLog(ctx context.Context, jobTime lgsv.JobTime, reportType lgsv.ReportType, recipients []string) {
	s.logger.Infof("Running reportLog type %d for %s", reportType, recipients)
	currentTime := time.Now()

	switch reportType {
	case lgsv.ReportType_SAVE_REPORT:
		startDate := s.getStartDate(currentTime.Unix(), jobTime)
		logs := s.logRepository.GetLog(ctx, startDate, currentTime.Unix(), nil, nil, "")
		csvData := [][]string{
			{"Company", "User ID", "User Name", "DeviceUuid", "Activity", "Date"},
		}
		for _, log := range logs {
			s.logger.Infof("LOG: %s in %s do %d at %s", log.UserId, log.DeviceUuid, log.Activity, time.Unix(log.Date, 0).Local().Format(time.RFC3339))
			csvData = append(csvData, []string{log.CompanyCode, log.UserId, log.UserName, log.DeviceUuid, strconv.Itoa(int(log.Activity)), time.Unix(log.Date, 0).Local().Format(time.RFC3339)})
		}
		s.logger.Infof("Save to csv file")
		byteBuf := writeToCSV(csvData)
		file_name := fmt.Sprintf("/tmp/Report_From_%s_To_%s.csv", time.Unix(startDate, 0).Local().Format(time.RFC3339), currentTime.Local().Format(time.RFC3339))
		f, _ := os.Create(file_name)
		defer f.Close()
		n2, _ := f.Write(byteBuf)
		s.logger.Infof("Wrote %d bytes\n", n2)
		f.Sync()
		break
	case lgsv.ReportType_SEND_EMAIL:
		break
	default:
		s.logger.Infof("Cannot report")
	}
}

func (s *service) getStartDate(currentTime int64, jobTime lgsv.JobTime) int64 {
	startDate := int64(0)

	switch jobTime {
	case lgsv.JobTime_EVERY_HOUR:
		startDate = currentTime - 3600
		break
	case lgsv.JobTime_EVERY_DAY:
		startDate = currentTime - 24*3600
		break
	case lgsv.JobTime_EVERY_MONTH:
		startDate = currentTime - 31*24*3600
		break
	case lgsv.JobTime_EVERY_WEEK:
	case lgsv.JobTime_EVERY_SUNDAY:
	case lgsv.JobTime_EVERY_MONDAY:
	case lgsv.JobTime_EVERY_TUESDAY:
	case lgsv.JobTime_EVERY_WEDNESDAY:
	case lgsv.JobTime_EVERY_THURSDAY:
	case lgsv.JobTime_EVERY_FRIDAY:
	case lgsv.JobTime_EVERY_SATURDAY:
		startDate = currentTime - 7*24*3600
		break
	default:
		s.logger.Infof("Cannot schedule report")
	}
	return startDate
}

func writeToCSV(empData [][]string) []byte {
	var buffer bytes.Buffer

	csvwriter := csv.NewWriter(&buffer)

	for _, empRow := range empData {
		_ = csvwriter.Write(empRow)
	}

	csvwriter.Flush()
	return buffer.Bytes()
}

func (s *service) handleAuditLogConsumer() {

	err := s.auditLogConsumer.SubscribeTopics([]string{constants.KAFKATOPIC_DATA_UPDATED}, nil)
	if err != nil {
		s.logger.Errorf("Error at SubscribeTopics KAFKATOPIC_DATA_UPDATED: %v", err)
		return
	}
	topics, err := s.auditLogConsumer.Subscription()
	s.logger.Infof("Subscribed to: %v\n", topics)
	if err != nil {
		s.logger.Errorf("Error at Subscription KAFKATOPIC_DATA_UPDATED: %v", err)
		return
	}

	for {
		ev := s.auditLogConsumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			msg := strings.TrimSpace(string(e.Value))
			s.logger.Infof("%% Message on topic %s\n -- msg: %v\n", *e.TopicPartition.Topic, msg)
			if msg != "" {
				s.logRepository.SaveAuditLog(msg, s.getMongodbConfig())
			}

		case kafka.PartitionEOF:
			s.logger.Infof("%% Reached %v\n", e)
		case kafka.Error:
			s.logger.Errorf("%% Error: %v\n", e)
			continue
		default:
		}
	}

}

func (s *service) handleVerifyConsumer() {
	err := s.kafkaConsumer.SubscribeTopics([]string{USER_SERVICE_VERIFICATION_TOPIC}, nil)
	if err != nil {
		s.logger.Errorf("handleVerifyConsumer unable to subscribe topic USER_SERVICE_VERIFICATION_TOPIC: %s\n", err)
		return
	}
	// Process messages
	sigchan := make(chan os.Signal, 1)

	for {
		select {
		case sig := <-sigchan:
			s.logger.Infof("Caught signal %v: terminating\n", sig)
			s.kafkaConsumer.Close()
			return
		default:
			msg, err := s.kafkaConsumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				continue
			}
			recordValue := msg.Value
			var user_event entities.UserVerificationEvent
			err = json.Unmarshal(recordValue, &user_event)
			if err == nil {
				s.logRepository.SaveLog(user_event.AppUserId, constants.VERIFY, user_event.DeviceId, string(recordValue), user_event.VerificationTime)
			} else {
				s.logger.Infof("Unmarshal error : %s\n", err)
			}
		}
	}
}

func (s *service) getMongodbConfig() viewmodels.MongoDbConfig {
	url := "mongodb://" +
		s.cfg.Mongodb.MongodbUser +
		":" + s.cfg.Mongodb.MongodbPassword +
		"@" + s.cfg.Mongodb.MongodbHost +
		":" + s.cfg.Mongodb.MongodbPort +
		"/" + s.cfg.Mongodb.MongodbDbname +
		"?authSource=" + s.cfg.Mongodb.MongodbDbname
	s.logger.Infof("getMongodbConfig: %v", url)
	mongoConfig := viewmodels.MongoDbConfig{
		UserName:   s.cfg.Mongodb.MongodbUser,
		Password:   s.cfg.Mongodb.MongodbPassword,
		Url:        url,
		Database:   s.cfg.Mongodb.MongodbDbname,
		Collection: "AuditLog",
	}
	return mongoConfig
}
