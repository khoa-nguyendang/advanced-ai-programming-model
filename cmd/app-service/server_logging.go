package main

import (
	"aapi/pkg/minio_srv"
	appsvc "aapi/protos/v1"
	lgsv "aapi/services/logging/protos/log/v1"
	"aapi/shared/constants"
	httpclients "aapi/shared/http_clients"
	vm "aapi/shared/view_models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var company_code string = ""
var deviceUUid string = ""

func (s *appServer) GetLog(ctx context.Context, model *appsvc.QueryRequest) (*appsvc.QueryResponse, error) {
	if model == nil {
		s.logger.Info("Get request is empty")
		return nil, errors.New("Get request is empty")
	}

	s.logger.Infof("GetLog StartDate=%v --- EndDate=%v --- DeviceUuids=%v --- UserIds=%v --- SaveToFile=%v \n", model.StartDate, model.EndDate, model.DeviceUuids, model.UserIds, model.SaveToFile)
	data, err := s.loggingApiClient.GetLog(ctx, &lgsv.QueryRequest{
		StartDate:   model.StartDate,
		EndDate:     model.EndDate,
		DeviceUuids: model.DeviceUuids,
		UserIds:     model.UserIds,
		SaveToFile:  model.SaveToFile,
	})

	if err != nil {
		s.logger.Errorf("loggingApiClient.GetLog: %v", err)
		return nil, err
	}

	if data == nil {
		s.logger.Errorf("loggingApiClient.GetLog empty response.")
		return &appsvc.QueryResponse{
			Logs:       make([]*appsvc.Log, 0),
			ReportFile: make([]byte, 0),
		}, nil
	}

	var logs []*appsvc.Log
	for _, log := range data.Logs {
		logs = append(logs, &appsvc.Log{
			Id:          log.Id,
			DeviceUuid:  log.DeviceUuid,
			Activity:    log.Activity,
			Date:        log.Date,
			UserId:      log.UserId,
			CompanyCode: log.CompanyCode,
			UserName:    log.UserName,
			UserState:   appsvc.UserState(log.UserState),
		})
	}

	return &appsvc.QueryResponse{
		Logs:       logs,
		ReportFile: data.ReportFile,
	}, err
}

func (s *appServer) GetLogs(ctx context.Context, model *appsvc.GetLogsRequest) (*appsvc.GetLogsResponse, error) {
	if model == nil {
		s.logger.Info("Get request is empty")
		return nil, errors.New("Get request is empty")
	}

	data, err := s.loggingApiClient.GetLogs(ctx, &lgsv.GetLogsRequest{
		StartDate:   model.StartDate,
		EndDate:     model.EndDate,
		DeviceUuids: model.DeviceUuids,
		UserIds:     model.UserIds,
		SaveToFile:  model.SaveToFile,
		Page:        model.Page,
		Size:        model.Size,
		CompanyCode: model.CompanyCode,
	})

	if err != nil || data == nil {
		s.logger.Errorf("loggingApiClient.GetLog: %v", err)
		return nil, err
	}

	var logs []*appsvc.Log
	for _, log := range data.Data {
		logs = append(logs, &appsvc.Log{
			Id:          log.Id,
			DeviceUuid:  log.DeviceUuid,
			Activity:    log.Activity,
			Date:        log.Date,
			UserId:      log.UserId,
			CompanyCode: log.CompanyCode,
			UserName:    log.UserName,
			UserState:   appsvc.UserState(log.UserState),
		})
	}

	return &appsvc.GetLogsResponse{
		Data:       logs,
		ReportFile: data.ReportFile,
		Total:      data.Total,
	}, err
}

func (s *appServer) ScheduleReport(ctx context.Context, model *appsvc.ScheduleReportRequest) (*appsvc.ScheduleReportResponse, error) {
	if model == nil {
		s.logger.Info("Get request is empty")
		return nil, errors.New("Get request is empty")
	}

	data, err := s.loggingApiClient.ScheduleReport(ctx, &lgsv.ScheduleReportRequest{
		JobTime:      lgsv.JobTime(model.JobTime),
		JobTimeCount: model.JobTimeCount,
		ReportType:   lgsv.ReportType(model.ReportType),
		Recipients:   model.Recipients,
	})

	if err != nil || data == nil {
		s.logger.Errorf("loggingApiClient.ScheduleReport: %v", err)
		return nil, err
	}

	return &appsvc.ScheduleReportResponse{
		Data: data.Data,
	}, err
}

func (s *appServer) AddLog(ctx context.Context, model *appsvc.AddLogRequest) (*appsvc.AddLogResponse, error) {
	if model == nil {
		s.logger.Info("AddLog request is empty")
		return nil, errors.New("AddLog request is empty")
	}

	var data *lgsv.AddLogResponse
	var err error
	var mInput *vm.MatchingInput
	if model.Title == "unknown" {
		mInput = s.handleVerificationFailed(ctx, model.DeviceUuid, model.CompanyCode, model.Image)
	} else {
		mInput = s.handleVerificationSuccess(ctx, model.DeviceUuid, model.CompanyCode, model.Image, model.Body)
	}

	if mInput == nil {
		return &appsvc.AddLogResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse model",
			Error:   "FAILED_TO_PARSE_MODEL",
			Status:  false,
		}, nil
	}

	newJsonBody, err := json.Marshal(mInput)

	if err != nil {
		s.logger.Errorf("AddLog.MatchingInput marshall error: %v", err)
		return &appsvc.AddLogResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   "FAILED_TO_PARSE_MODEL",
			Status:  false,
		}, nil
	}

	model.Body = string(newJsonBody)
	data, err = s.loggingApiClient.AddLog(ctx, &lgsv.AddLogRequest{
		Title:       model.Title,
		Body:        model.Body,
		DeviceUuid:  model.DeviceUuid,
		CompanyCode: model.CompanyCode,
		Activity:    model.Activity,
	})

	if err != nil {
		s.logger.Errorf("loggingApiClient.AddLog.got error: %v", err)
		return &appsvc.AddLogResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   err.Error(),
			Status:  false,
		}, nil
	}

	if data == nil {
		s.logger.Errorf("loggingApiClient.AddLog.got data nil: %v", err)
		return &appsvc.AddLogResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: "AddLog service return empty",
			Error:   "RETURN_NULL_DATA",
			Status:  false,
		}, nil
	}

	return &appsvc.AddLogResponse{
		Code:    http.StatusOK,
		Message: data.Message,
		Error:   data.Error,
		Status:  data.Status,
	}, err
}

func (s *appServer) broadcastUserVerify(wg *sync.WaitGroup, verifyJson, deviceUUid, company_code string) {
	defer wg.Done()
	message := vm.CommonMessage{
		MessageType:  constants.VERIFICATION_EVENT,
		From:         deviceUUid,
		To:           "WebAdmin",
		Topic:        company_code,
		CompanyCode:  company_code,
		Body:         verifyJson,
		CreatedDate:  time.Now().UTC().UnixMilli(),
		LastModified: time.Now().UTC().UnixMilli(),
		Acknowledge:  false,
	}
	s.logger.Infof("broadcastUserVerify: %v", message)
	url := s.cfg.PublisherServer + "/send-message"
	responseMessage, err := httpclients.SendHttpPost2(url, constants.Annonymous, message, "", "", "")
	if err != nil {
		s.logger.Errorf("broadcastUserVerify got error: %v", err)
	}

	s.logger.Infof("broadcastUserVerify.success: %v", string(responseMessage))
}

// handleVerificationFailed
func (s *appServer) handleVerificationFailed(ctx context.Context, deviceUUid, company_code string, img *appsvc.Image) *vm.MatchingInput {
	imgPath, preSignedUrl := s.uploadImageToMinio(img, company_code, false)

	dt := &vm.MatchingInput{
		UserId:       "0",
		UserName:     "",
		Temperature:  "N/A",
		DeviceName:   deviceUUid,
		ImagePathS3:  imgPath,
		PreSignedUrl: preSignedUrl,
		VerifyStatus: false,
	}

	jsondata, err := json.Marshal(dt)

	if err != nil {
		s.logger.Errorf("handleVerificationFailed Unable to marshal object")
		return nil
	}

	s.handleVerification(ctx, string(jsondata), deviceUUid, company_code)
	return dt
}

func (s *appServer) handleVerificationSuccess(ctx context.Context,
	deviceUUid, company_code string,
	img *appsvc.Image,
	body string) *vm.MatchingInput {

	var dt *vm.MatchingInput = &vm.MatchingInput{}
	if body == "" {
		dt = &vm.MatchingInput{
			UserId:       "0",
			UserName:     "",
			Temperature:  "N/A",
			DeviceName:   "",
			ImagePathS3:  "",
			PreSignedUrl: "",
			VerifyStatus: false,
			UserStatus:   false,
		}
	} else {
		err := json.Unmarshal([]byte(body), dt)
		if err != nil {
			s.logger.Errorf("Unable to Unmarshal body AddLog request %v \n, err: %v", body, err)
			return nil
		}

	}

	imgPath, preSignedUrl := s.uploadImageToMinio(img, company_code, true)
	dt.ImagePathS3 = imgPath
	dt.PreSignedUrl = preSignedUrl
	dt.VerifyStatus = true

	matchingJson, err := json.Marshal(dt)
	if err != nil {
		s.logger.Errorf("handleVerificationSuccess got error when marshal MatchingInput struct: %v", err)
		return dt
	}
	s.handleVerification(ctx, string(matchingJson), deviceUUid, company_code)
	return dt
}

func (s *appServer) handleVerification(ctx context.Context, userJson, deviceUUid, company_code string) {
	message := vm.CommonMessage{
		MessageType:  constants.VERIFICATION_EVENT,
		From:         "Device" + deviceUUid,
		To:           "WebAdmin",
		Topic:        company_code,
		Body:         userJson,
		CompanyCode:  company_code,
		CreatedDate:  time.Now().UTC().UnixMilli(),
		LastModified: time.Now().UTC().UnixMilli(),
		Acknowledge:  false,
	}
	s.broadcastVerification(message)
}

func (s *appServer) broadcastVerification(message vm.CommonMessage) {
	s.logger.Infof("broadcastVerification: %v", message)
	url := s.cfg.PublisherServer + "/send-message"
	responseMessage, err := httpclients.SendHttpPost2(url, constants.Annonymous, message, "", "", "")
	if err != nil {
		s.logger.Errorf("broadcastVerification got error: %v", err)
	}

	s.logger.Infof("broadcastVerification.success: %v", string(responseMessage))
}

// uploadImageToMinio upload a picture to minio then return path, presigned url
func (s *appServer) uploadImageToMinio(img *appsvc.Image, company_code string, status bool) (path string, presignedUrl string) {
	if img == nil {
		return "", ""
	}
	var statusPath string

	if status {
		statusPath = "/success/"
	} else {
		statusPath = "/fail/"
	}

	if company_code == "" {
		company_code = "default"
	}
	currentDate := time.Now().UTC().Format("2006-02-01")
	imgId := strconv.FormatInt(time.Now().UTC().UnixMicro(), 10)
	target_file_name := "verification/" + currentDate + statusPath + imgId + ".jpg"
	file_reader := bytes.NewReader(img.Data)
	return target_file_name, minio_srv.UploadFile(s.minioClient, constants.MINIOBUCKET_UserVerification, target_file_name, file_reader, file_reader.Size())
}
