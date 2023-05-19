package repositories

import (
	"aapi/pkg/mongodb"
	lsv "aapi/services/logging/protos/log/v1"
	"aapi/shared/constants"
	dbqueries "aapi/shared/database_queries"
	"aapi/shared/entities"
	viewmodels "aapi/shared/view_models"
	"context"
	"database/sql"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func (rep *logRepo) GetLog(ctx context.Context, start_date int64, end_date int64, device_uuids []string, user_ids []string, company_code string) []*entities.Log {
	var logs []*entities.Log

	var rows *sql.Rows
	var err error
	if company_code != "" {
		if device_uuids != nil && len(device_uuids) > 0 {
			if user_ids != nil && len(user_ids) > 0 {
				query, args, err := sqlx.In(dbqueries.FindLogWithCompanyCodeQueryByUserIdAndDeviceUuid, company_code, start_date, end_date, user_ids, device_uuids)
				if err != nil {
					rep.logger.Infof("Error %v", err)
					return logs
				}
				query = rep.db.Rebind(query)
				rows, err = rep.db.Query(query, args...)
			} else {
				rep.logger.Infof("Device UUID %v", device_uuids[0])
				query, args, err := sqlx.In(dbqueries.FindLogWithCompanyCodeQueryByDeviceUuid, company_code, start_date, end_date, device_uuids)
				if err != nil {
					rep.logger.Infof("Error %v", err)
					return logs
				}
				query = rep.db.Rebind(query)
				rows, err = rep.db.Query(query, args...)
			}
			if err != nil {
				rep.logger.Infof("Error %v", err)
				return logs
			}
		} else {
			if user_ids != nil && len(user_ids) > 0 {
				rep.logger.Infof("user_id %v", user_ids[0])

				query, args, err := sqlx.In(dbqueries.FindLogWithCompanyCodeQueryByUserId, company_code, start_date, end_date, user_ids)
				if err != nil {
					return logs
				}
				query = rep.db.Rebind(query)
				rows, err = rep.db.Query(query, args...)
			} else {
				rows, err = rep.db.Query(dbqueries.FindLogWithCompanyCodeQuery, company_code, start_date, end_date)
			}
			if err != nil {
				rep.logger.Infof("Error %v", err)
				return logs
			}
		}
	} else {
		if device_uuids != nil && len(device_uuids) > 0 {
			if user_ids != nil && len(user_ids) > 0 {
				query, args, err := sqlx.In(dbqueries.FindLogQueryByUserIdAndDeviceUuid, start_date, end_date, user_ids, device_uuids)
				if err != nil {
					rep.logger.Infof("Error %v", err)
					return logs
				}
				query = rep.db.Rebind(query)
				rows, err = rep.db.Query(query, args...)
			} else {
				rep.logger.Infof("Device UUID %v", device_uuids[0])
				query, args, err := sqlx.In(dbqueries.FindLogQueryByDeviceUuid, start_date, end_date, device_uuids)
				if err != nil {
					rep.logger.Infof("Error %v", err)
					return logs
				}
				query = rep.db.Rebind(query)
				rows, err = rep.db.Query(query, args...)
			}
			if err != nil {
				rep.logger.Infof("Error %v", err)
				return logs
			}
		} else {
			if user_ids != nil && len(user_ids) > 0 {
				rep.logger.Infof("user_id %v", user_ids[0])

				query, args, err := sqlx.In(dbqueries.FindLogQueryByUserId, start_date, end_date, user_ids)
				if err != nil {
					return logs
				}
				query = rep.db.Rebind(query)
				rows, err = rep.db.Query(query, args...)
			} else {
				rows, err = rep.db.Query(dbqueries.FindLogQuery, start_date, end_date)
			}
			if err != nil {
				rep.logger.Infof("Error %v", err)
				return logs
			}
		}
	}

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			log := &entities.Log{}

			err = rows.Scan(&log.Id,
				&log.AppUserId,
				&log.Activity,
				&log.DeviceUuid,
				&log.FullMessage,
				&log.Date,
				&log.CompanyCode,
				&log.UserId,
				&log.UserName,
				&log.UserState,
			)

			if err == nil {
				logs = append(logs, log)
			} else {
				rep.logger.Infof("Error %v", err)
			}
		}
	} else {
		rep.logger.Infof("Error %v", err)
	}

	return logs
}

func (rep *logRepo) GetLogs(ctx context.Context, model *lsv.GetLogsRequest) (*lsv.GetLogsResponse, error) {
	rep.logger.Infof("GetLogs: %v", model.CompanyCode)
	var companyCode string = model.CompanyCode
	var logs []*lsv.Log = make([]*lsv.Log, 0)
	var deviceId string = ""
	var userId string = ""
	var total int32 = 0
	var fromDateTime int64 = 0
	var toDateTime int64 = 9999999999
	var rows *sql.Rows
	var err error
	if len(model.DeviceUuids) > 0 {
		deviceId = model.DeviceUuids[0]
	}
	if len(model.UserIds) > 0 {
		userId = model.UserIds[0]
	}
	if model.StartDate > 0 {
		fromDateTime = model.StartDate
	}
	if model.EndDate > 0 {
		toDateTime = model.EndDate
	}

	if deviceId != "" && userId == "" {
		// case 1, device provided, but userId not
		rows, err = rep.db.Query(dbqueries.LOG_GetAllBy_CompanyCode_DeviceUUID, companyCode, deviceId, fromDateTime, toDateTime, model.Page*model.Size, model.Size)
		err = rep.db.GetContext(ctx, &total, dbqueries.LOG_CountGetAllBy_CompanyCode_DeviceUUID, companyCode, deviceId, fromDateTime, toDateTime)
	} else if deviceId != "" && userId != "" {
		// case 2, device , userId both are provided
		rows, err = rep.db.Query(dbqueries.LOG_GetAllBy_CompanyCode_DeviceUUID_UserId, companyCode, deviceId, userId, fromDateTime, toDateTime, model.Page*model.Size, model.Size)
		err = rep.db.GetContext(ctx, &total, dbqueries.LOG_CountGetAllBy_CompanyCode_DeviceUUID_UserId, companyCode, deviceId, userId, fromDateTime, toDateTime)
	} else if deviceId == "" && userId != "" {
		// case 3, userId provided, but device
		rows, err = rep.db.Query(dbqueries.LOG_GetAllBy_CompanyCode_UserId, companyCode, userId, fromDateTime, toDateTime, model.Page*model.Size, model.Size)
		err = rep.db.GetContext(ctx, &total, dbqueries.LOG_CountGetAllBy_CompanyCode_UserId, companyCode, userId, fromDateTime, toDateTime)
	} else if deviceId == "" && userId == "" {
		rep.logger.Info("deviceId and userId is empty")
		// case 4, neither device, userId was provided
		rows, err = rep.db.Query(dbqueries.LOG_GetAllBy_CompanyCode, companyCode, fromDateTime, toDateTime, model.Page*model.Size, model.Size)
		err = rep.db.GetContext(ctx, &total, dbqueries.LOG_CountGetAllBy_CompanyCode, companyCode, fromDateTime, toDateTime)
	}

	if err != nil {
		rep.logger.Errorf("err: %v", err)
		return &lsv.GetLogsResponse{
			Data:       make([]*lsv.Log, 0),
			ReportFile: make([]byte, 0),
			Total:      0,
			Code:       http.StatusBadRequest,
			Message:    err.Error(),
			Error:      "FAILED_TO_GET_LOGS",
		}, nil
	}
	i := 0
	defer rows.Close()
	for rows.Next() {
		i++
		log := &entities.Log{}
		err = rows.Scan(&log.Id,
			&log.Activity,
			&log.DeviceUuid,
			&log.Date,
			&log.CompanyCode,
			&log.UserId,
			&log.UserName,
			&log.UserState,
		)

		if err == nil {
			logs = append(logs, &lsv.Log{
				CompanyCode: log.CompanyCode,
				UserId:      log.UserId,
				UserName:    log.UserName,
				UserState:   lsv.UserState(log.UserState),
				Activity:    int64(log.Activity),
				DeviceUuid:  log.DeviceUuid,
				Date:        log.Date,
				Id:          log.Id,
			})
		} else {
			rep.logger.Infof("Error %v", err)
		}
	}
	rep.logger.Infof("i: %v", i)

	return &lsv.GetLogsResponse{
		Data:       logs,
		ReportFile: make([]byte, 0),
		Total:      total,
		Code:       http.StatusOK,
		Message:    "",
		Error:      "",
	}, err
}

func (rep *logRepo) SaveLog(app_user_id int64, activity constants.Activity, device_uuid string, full_message string, date int64) error {
	ctx := context.Background()
	res, err := rep.db.ExecContext(ctx,
		dbqueries.SaveLogQuery,
		app_user_id,
		activity,
		device_uuid,
		full_message,
		date,
	)
	if err != nil || res == nil {
		rep.logger.Infof("SaveLog.err: %v", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected <= 0 {
		rep.logger.Infof("SaveLog.err or no rows affects: %v", err)
		return err
	}
	return nil
}

func (rep *logRepo) SaveAuditLog(auditLogJson string, mdc viewmodels.MongoDbConfig) error {
	err := mongodb.InsertJson(rep.mongoClient, mdc, constants.MGC_AUDIT_LOG, auditLogJson)
	rep.logger.Infof("SaveAuditLog: %#v---, %v", mdc, err)
	return err
}
