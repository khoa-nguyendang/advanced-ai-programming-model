package repositories

import (
	user_sv "aapi/services/user/protos/user/v1"
	"aapi/shared/constants"
	dbqueries "aapi/shared/database_queries"
	"aapi/shared/entities"
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func (u *userRepo) ValidateData(ctx context.Context, company_code string, group_ids []int64) (int32, string, string) {
	return 0, "", ""

	// Check company code
	var existing_company_id int64 = 1
	var existing_company_code string
	if err := u.db.QueryRowContext(ctx, dbqueries.CheckCompanyCodeExistQuery, company_code).Scan(&existing_company_id, &existing_company_code); err != nil && errors.Cause(err) != sql.ErrNoRows {
		u.logger.Errorf("CheckCompanyCodeExistQuery err : %v", err)
		return http.StatusInternalServerError, "Failed to query database", "INTERNAL_SERVER_ERROR"
	}

	u.logger.Infof("Found company code: %v", existing_company_code)

	if existing_company_code == "" {
		// company code doesn't exist.
		return http.StatusBadRequest, "Company code does not exist. Please contact admin", "COMPANY_CODE_DOES_NOT_EXIST"
	}

	groups, _ := u.GetGroupsByCompanyId(ctx, existing_company_id)
	for _, group_id := range group_ids {
		if group_id != 0 {
			is_group_valid := false
			for _, group := range groups {
				if group_id == group.Id {
					is_group_valid = true
					break
				}
			}
			if !is_group_valid {
				return http.StatusBadRequest, "User group is not valid", "USER_GROUP_IS_NOT_VALID"
			}
		}
	}

	return 0, "", ""
}

func (u *userRepo) EnrollUser(ctx context.Context, user *entities.User, image_URLs []string) (*entities.User, error) {
	u.logger.Infof("EnrollUser UserId: %v", user.UserId)

	current_time := time.Now().UTC().Unix()

	result, err := u.db.ExecContext(
		ctx,
		dbqueries.EnrollUserQuery,
		user.CompanyCode,
		user.UserId,
		user.UserName,
		user.UserRoleId,
		user.UserInfo,
		user.UserState,
		user.ThumbnailImageUrl,
		current_time,
		current_time,
		user.ActivationDate,
		user.ExpiryDate,
		"TBD",
	)
	u.logger.Infof("EnrollUser execute database: %v", result)
	if err != nil {
		u.logger.Errorf("EnrollUser execute database.err: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	for _, image_path := range image_URLs {
		u.logger.Infof("EnrollUser save image %v for user: %v", image_path, id)
		result, err = u.db.ExecContext(
			ctx,
			dbqueries.InsertEnollmentImage,
			id,
			image_path,
			0,
			current_time,
			current_time,
			0,
			0,
		)
		if err != nil {
			u.logger.Errorf("EnrollUser execute database.err: %v", err)
		}
	}

	for _, group_id := range user.UserGroupIds {
		u.logger.Infof("EnrollUser save group %v for user: %v", group_id, id)
		result, err = u.db.ExecContext(
			ctx,
			dbqueries.InsertUserGroups,
			id,
			group_id,
			1,
			current_time,
			current_time,
			0,
			0,
		)
		if err != nil {
			u.logger.Errorf("EnrollUser execute database.err: %v", err)
		}
	}

	u.logger.Infof("EnrollUser success: %v", user)

	user.Id = id
	user.CreatedAt = current_time
	user.LastModified = current_time
	return user, nil
}

func (u *userRepo) FindUserByUserId(ctx context.Context, company_code string, user_id string) (*entities.User, error) {
	user := &entities.User{}
	if company_code == "" {
		if err := u.db.QueryRowContext(ctx, dbqueries.FindUserByUserIDQuery, user_id, constants.USER_STATE_DELETED).Scan(&user.Id, &user.CompanyCode,
			&user.UserId,
			&user.UserName,
			&user.UserRoleId,
			&user.UserInfo,
			&user.UserState,
			&user.ThumbnailImageUrl,
			&user.LastModified,
			&user.IssuedDate,
			&user.ActivationDate,
			&user.ExpiryDate,
			&user.ReferenceId); err != nil {
			u.logger.Infof("FindUser error for: %v, err: %v", user_id, err)
			if errors.Cause(err) == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
	} else {
		if err := u.db.QueryRowContext(ctx, dbqueries.FindUserWithCompanyCodeByUserIDQuery, company_code, user_id, constants.USER_STATE_DELETED).Scan(&user.Id, &user.CompanyCode,
			&user.UserId,
			&user.UserName,
			&user.UserRoleId,
			&user.UserInfo,
			&user.UserState,
			&user.ThumbnailImageUrl,
			&user.LastModified,
			&user.IssuedDate,
			&user.ActivationDate,
			&user.ExpiryDate,
			&user.ReferenceId); err != nil {
			u.logger.Infof("FindUser error for: %v, err: %v", user_id, err)
			if errors.Cause(err) == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
	}
	u.logger.Infof("FindUser success for user: %v, data: %v", user_id, user)
	return user, nil
}

func (u *userRepo) FindAllUser(ctx context.Context) []*entities.User {
	var users []*entities.User
	rows, err := u.db.Query(dbqueries.FindAllUserQuery, constants.USER_STATE_DELETED)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			user := &entities.User{}

			err = rows.Scan()
			u.logger.Error(err)

			if err == nil {
				users = append(users, user)
			}
		}
	}
	return users
}

func (u *userRepo) GetUser(ctx context.Context, model *user_sv.GetUserRequest, company_code string) (*user_sv.GetUserResponse, error) {
	users := []*user_sv.UserData{}
	entity_users := []*entities.User{}

	var rows *sql.Rows
	var err error
	var total_count int64

	limit := int64(1000)
	skip := int64(0)
	if model.PageSize > 0 {
		limit = model.PageSize
	}
	if model.CurrentPage > 0 {
		skip = limit * model.CurrentPage
	}

	u.logger.Infof("GetUser company_code = %v UserRoleId = %v skip = %v limit = %v ", company_code, model.UserRoleId, skip, limit)

	err = u.db.QueryRowContext(ctx, dbqueries.CountUserQuery, constants.USER_STATE_DELETED).Scan(&total_count)
	u.logger.Infof("Error %v total_count: %v", err, total_count)
	rows, err = u.db.Query(dbqueries.FindAllUserQuery, constants.USER_STATE_DELETED, skip, limit)

	user_column_ids := []int64{}
	if err == nil && rows != nil {
		defer rows.Close()
		for rows.Next() {
			user := &entities.User{}
			err = rows.Scan(
				&user.Id,
				&user.CompanyCode,
				&user.UserId,
				&user.UserName,
				&user.UserRoleId,
				&user.UserInfo,
				&user.UserState,
				&user.ThumbnailImageUrl,
				&user.LastModified,
				&user.IssuedDate,
				&user.ActivationDate,
				&user.ExpiryDate,
				&user.ReferenceId,
			)

			if err == nil {
				entity_users = append(entity_users, user)
				user_column_ids = append(user_column_ids, user.Id)
			} else {
				u.logger.Infof("Error %v", err)
			}
		}
	} else {
		u.logger.Infof("Error %v", err)
	}

	user_image_paths := u.GetEnrollmentImages(ctx, user_column_ids)
	u.logger.Infof("Get enrollment image len = %v:", len(user_image_paths))

	for i := range entity_users {
		u.logger.Infof("Processing user at column id %v:", entity_users[i].Id)
		users = append(users, &user_sv.UserData{
			UserId:              entity_users[i].UserId,
			UserName:            entity_users[i].UserName,
			UserRoleId:          entity_users[i].UserRoleId,
			UserRole:            "role_name",
			UserGroupIds:        []int64{},
			UserGroups:          []string{},
			UserInfo:            entity_users[i].UserInfo,
			State:               user_sv.UserState(entity_users[i].UserState),
			ThumbnailImageUrl:   entity_users[i].ThumbnailImageUrl,
			RegisteredImageUrls: user_image_paths[entity_users[i].Id],
			LastModified:        entity_users[i].LastModified,
			IssuedDate:          entity_users[i].IssuedDate,
			ActivationDate:      entity_users[i].ActivationDate,
			ExpiryDate:          entity_users[i].ExpiryDate,
			IsActive:            true,
		})
	}

	return &user_sv.GetUserResponse{
		Users:       users,
		CurrentPage: model.CurrentPage,
		PageSize:    limit,
		TotalCount:  total_count,
	}, nil
}

func (u *userRepo) UpdateUser(ctx context.Context, user *entities.User, image_URLs []string, company_code string) (*entities.User, error) {
	u.logger.Infof("Processing update user %v:", user)

	current_time := time.Now().UTC().Unix()

	if _, err := u.db.ExecContext(ctx, dbqueries.UpdateUserByUserIdQuery,
		user.UserName,
		user.UserRoleId,
		user.UserInfo,
		user.UserState,
		current_time,
		user.ActivationDate,
		user.ExpiryDate,
		user.ReferenceId,
		user.UserId,
		company_code,
		constants.USER_STATE_DELETED); err != nil {
		return nil, err
	}
	u.logger.Infof("UpdateUser success data: %v", user)

	for _, image_path := range image_URLs {
		_, err := u.db.ExecContext(
			ctx,
			dbqueries.InsertEnollmentImage,
			user.Id,
			image_path,
			0,
			current_time,
			current_time,
			0,
			0,
		)
		if err != nil {
			return nil, err
		}
		u.logger.Infof("Saved new image %v for user %v", image_path, user.Id)
	}

	if _, err := u.db.ExecContext(ctx, dbqueries.DeleteUserGroupsByUserColId, user.Id); err != nil {
		u.logger.Errorf("Error %v:", err)
	}

	for _, group_id := range user.UserGroupIds {
		u.logger.Info("UpdateUser save group %v for user: %v", group_id, user.Id)
		_, err := u.db.ExecContext(
			ctx,
			dbqueries.InsertUserGroups,
			user.Id,
			group_id,
			1,
			current_time,
			current_time,
			0,
			0,
		)
		if err != nil {
			u.logger.Errorf("UpdateUser execute database.err: %v", err)
		}
	}
	return user, nil
}

func (u *userRepo) DeleteUser(ctx context.Context, company_code string, user_id string) (bool, error) {
	if company_code == "" {
		if _, err := u.db.ExecContext(ctx, dbqueries.DeleteUserQuery, constants.USER_STATE_DELETED, time.Now().UTC().Unix(), user_id); err != nil {
			return false, err
		}
	} else {
		if _, err := u.db.ExecContext(ctx, dbqueries.DeleteUserWithCompanyCodesQuery, constants.USER_STATE_DELETED, time.Now().UTC().Unix(), company_code, user_id); err != nil {
			return false, err
		}
	}
	u.logger.Infof("DeleteUser success data: %v", user_id)
	return true, nil
}

func (u *userRepo) GetRoles(ctx context.Context) ([]*entities.Role, error) {
	var roles []*entities.Role
	err := u.db.SelectContext(ctx, &roles, dbqueries.GetRoles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (u *userRepo) GetGroupsByCompanyId(ctx context.Context, company_id int64) ([]*entities.Group, error) {
	var groups []*entities.Group
	err := u.db.SelectContext(ctx, &groups, dbqueries.GetGroupsByCompanyId, company_id)
	if err != nil {
		u.logger.Errorf("Error %v:", err)
	}
	return groups, nil
}

func (u *userRepo) GetGroups(ctx context.Context) ([]*entities.Group, error) {
	var groups []*entities.Group
	err := u.db.SelectContext(ctx, &groups, dbqueries.GetGroups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (u *userRepo) GetEnrollmentImages(ctx context.Context, user_column_ids []int64) map[int64][]string {
	dict := make(map[int64][]string)

	if len(user_column_ids) > 0 {
		query, args, err := sqlx.In(dbqueries.GetEnrollmentImagePath, user_column_ids)
		if err == nil {
			query = u.db.Rebind(query)
			rows, err := u.db.Query(query, args...)
			if err == nil && rows != nil {
				defer rows.Close()
				for rows.Next() {
					var user_column_id int64
					var image_path string

					err = rows.Scan(
						&user_column_id,
						&image_path,
					)

					if err == nil {
						dict[user_column_id] = append(dict[user_column_id], image_path)
					} else {
						u.logger.Errorf("Error %v:", err)
					}
				}
			} else {
				u.logger.Errorf("Error %v:", err)
			}
		}

	}
	u.logger.Infof("Enrollment image %v:", dict)
	return dict
}

func (u *userRepo) GetUserGroupsInfo(ctx context.Context, user_column_ids []int64) map[int64][]*entities.Group {
	dict := make(map[int64][]*entities.Group)

	if len(user_column_ids) > 0 {
		query, args, err := sqlx.In(dbqueries.GetUserGroupsByUserColIds, user_column_ids)
		if err == nil {
			query = u.db.Rebind(query)
			rows, err := u.db.Query(query, args...)
			if err == nil && rows != nil {
				defer rows.Close()
				for rows.Next() {
					var user_column_id int64
					group := entities.Group{}

					err = rows.Scan(
						&user_column_id,
						&group.Id,
						&group.Name,
					)

					if err == nil {
						u.logger.Infof("Add group info for user_column_id %v:", user_column_id)
						dict[user_column_id] = append(dict[user_column_id], &group)
					} else {
						u.logger.Errorf("Error %v:", err)
					}
				}
			} else {
				u.logger.Errorf("Error %v:", err)
			}
		}

	}
	return dict
}

func (u *userRepo) IsRoleValid(ctx context.Context, roleId int64) bool {
	return true
}

func (u *userRepo) CheckUserActivaton(ctx context.Context) (map[string][]string, map[string][]string) {
	activated_users := map[string][]string{}
	deactivated_users := map[string][]string{}

	now := time.Now()
	start_time_of_today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()
	start_time_of_tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local).Unix()

	u.logger.Infof("Check activation_date in %v - %v", start_time_of_today, start_time_of_tomorrow)
	u.logger.Infof("Activate user")
	rows, err := u.db.Query(dbqueries.FindUserByStateInActivationDuration, 5, constants.USER_STATE_ENABLED, start_time_of_today, start_time_of_tomorrow, start_time_of_tomorrow)

	if err == nil && rows != nil {
		defer rows.Close()

		// updated_state_user_ids := []int64{}
		for rows.Next() {
			user := &entities.User{}
			err = rows.Scan(
				&user.Id,
				&user.CompanyCode,
				&user.UserId,
				&user.UserState)
			if err == nil {
				u.logger.Infof("Going to activate user_id %v of company %v", user.UserId, user.CompanyCode)
				// updated_state_user_ids = append(updated_state_user_ids, user.Id)
				user_ids := activated_users[user.CompanyCode]
				if user_ids == nil {
					user_ids = make([]string, 0)
				}
				user_ids = append(user_ids, user.UserId)
				activated_users[user.CompanyCode] = user_ids
			} else {
				u.logger.Errorf("Error %v:", err)
			}
		}
	} else {
		u.logger.Errorf("Error %v:", err)
	}

	u.logger.Infof("Dectivate user")
	rows, err = u.db.Query(dbqueries.FindUserByStateOutActivationDuration, 5, constants.USER_STATE_ENABLED, start_time_of_tomorrow, start_time_of_today)

	if err == nil && rows != nil {
		defer rows.Close()

		// updated_state_user_ids := []int64{}
		for rows.Next() {
			user := &entities.User{}
			err = rows.Scan(
				&user.Id,
				&user.CompanyCode,
				&user.UserId,
				&user.UserState)
			if err == nil {
				u.logger.Infof("Going to deactivate user_id %v of company %v", user.UserId, user.CompanyCode)
				// updated_state_user_ids = append(updated_state_user_ids, user.Id)
				user_ids := deactivated_users[user.CompanyCode]
				if user_ids == nil {
					user_ids = make([]string, 0)
				}
				user_ids = append(user_ids, user.UserId)
				deactivated_users[user.CompanyCode] = user_ids
			} else {
				u.logger.Errorf("Error %v:", err)
			}
		}
	} else {
		u.logger.Errorf("Error %v:", err)
	}

	return activated_users, deactivated_users
}
