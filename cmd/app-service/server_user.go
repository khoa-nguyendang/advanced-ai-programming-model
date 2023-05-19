package main

import (
	appsvc "aapi/protos/v1"
	user_sv "aapi/services/user/protos/user/v1"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jinzhu/copier"
)

// User service
func (s *appServer) Enroll(ctx context.Context, model *appsvc.EnrollRequest) (*appsvc.EnrollResponse, error) {
	s.logger.Infof("server_user.go Enroll")
	if model == nil {
		s.logger.Info("Enroll empty model")
		return &appsvc.EnrollResponse{
			Code:    http.StatusBadRequest,
			Message: "Enroll with empty request",
			Error:   "ENROLL_WITH_EMPTY_REQUEST",
		}, nil
	}
	s.logger.Infof(fmt.Sprintf("Enroll.UserId: %v", model.UserId))
	s.logger.Infof(fmt.Sprintf("Enroll.UserInfo: %v", model.UserInfo))
	s.logger.Infof(fmt.Sprintf("Enroll.Images: %v", len(model.Images)))
	images := make([]*user_sv.Image, len(model.Images))
	for idx, image := range model.Images {
		if image == nil {
			s.logger.Infof("image is null")
			continue
		}

		images[idx] = &user_sv.Image{
			Data:    image.Data,
			ImageId: image.ImageId,
		}
	}

	data, err := s.userApiClient.Enroll(ctx, &user_sv.EnrollRequest{
		UserId:         model.UserId,
		UserName:       model.UserName,
		UserInfo:       model.UserInfo,
		UserRoleId:     model.UserRoleId,
		UserGroupIds:   model.UserGroupIds,
		ExpiryDate:     model.ExpiryDate,
		ActivationDate: model.ActivationDate,
		Images:         images,
	})

	if err != nil || data == nil {
		return &appsvc.EnrollResponse{
			Code:    http.StatusInternalServerError,
			Message: "Extract feature error",
			Error:   "INTERNAL_SERVER_ERROR",
		}, nil
	}
	return &appsvc.EnrollResponse{
		Code:    data.Code,
		Error:   data.Error,
		Message: data.Message,
		State:   appsvc.UserState(data.State),
	}, nil
}

func (s *appServer) Verify(ctx context.Context, model *appsvc.VerifyRequest) (*appsvc.VerifyResponse, error) {

	if model == nil {
		s.logger.Info("Verify empty model")
		return &appsvc.VerifyResponse{
			Code:    http.StatusBadRequest,
			Message: "Verify with empty request",
			Error:   "VERIFY_WITH_EMPTY_REQUEST",
		}, nil
	}
	s.logger.Infof("server_user.go Verify")
	var images []*user_sv.Image

	if model.Images != nil {
		for _, image := range model.Images {
			if image != nil && image.Data != nil {
				images = append(images, &user_sv.Image{
					Data:    image.Data,
					ImageId: image.ImageId,
				})
			}
		}
	}

	data, err := s.userApiClient.Verify(ctx, &user_sv.VerifyRequest{
		Images: images,
	})

	if err != nil || data == nil {
		return &appsvc.VerifyResponse{
			Code:    http.StatusInternalServerError,
			Message: "Extract feature error",
			Error:   "INTERNAL_SERVER_ERROR",
		}, nil
	}

	return &appsvc.VerifyResponse{
		Code:              data.Code,
		Error:             data.Error,
		Message:           data.Message,
		UserId:            data.UserId,
		UserName:          data.UserName,
		State:             appsvc.UserState(data.State),
		UserRole:          data.UserRole,
		UserGroups:        data.UserGroups,
		UserInfo:          data.UserInfo,
		ThumbnailImageUrl: data.ThumbnailImageUrl,
		ImageId:           data.ImageId,
		Score:             data.Score,
		LastModified:      data.LastModified,
		IssuedDate:        data.IssuedDate,
		ExpiryDate:        data.ExpiryDate,
		ActivationDate:    data.ActivationDate,
		IsActive:          data.IsActive,
	}, nil
}

func (s *appServer) GetUser(ctx context.Context, model *appsvc.GetUserRequest) (*appsvc.GetUserResponse, error) {
	if model == nil {
		s.logger.Info("Get request is empty")
		return nil, errors.New("GET_REQUEST_IS_EMPTY")
	}

	data, err := s.userApiClient.GetUser(ctx, &user_sv.GetUserRequest{
		UserRoleId:   model.UserRoleId,
		UserGroupIds: model.UserGroupIds,
		UserIds:      model.UserIds,
		CurrentPage:  model.CurrentPage,
		PageSize:     model.PageSize,
	})

	if err != nil || data == nil {
		return nil, err
	}

	var users []*appsvc.UserData
	for _, user := range data.Users {
		users = append(users, &appsvc.UserData{
			UserId:              user.UserId,
			UserName:            user.UserName,
			UserRole:            user.UserRole,
			UserGroups:          user.UserGroups,
			UserInfo:            user.UserInfo,
			LastModified:        user.LastModified,
			UserGroupIds:        user.UserGroupIds,
			RegisteredImageUrls: user.RegisteredImageUrls,
			ThumbnailImageUrl:   user.ThumbnailImageUrl,
			State:               appsvc.UserState(user.State),
			UserRoleId:          user.UserRoleId,
			IssuedDate:          user.IssuedDate,
			ExpiryDate:          user.ExpiryDate,
			ActivationDate:      user.ActivationDate,
			IsActive:            user.IsActive,
		})
	}

	return &appsvc.GetUserResponse{
		Users:       users,
		CurrentPage: data.CurrentPage,
		PageSize:    data.PageSize,
		TotalCount:  data.TotalCount,
	}, err
}

func (s *appServer) CountUser(ctx context.Context, model *appsvc.CountUserRequest) (*appsvc.CountUserResponse, error) {
	if model == nil {
		s.logger.Info("CountUser request is empty")
		return nil, errors.New("COUNT_REQUEST_IS_EMPTY")
	}

	data, err := s.userApiClient.CountUser(ctx, &user_sv.CountUserRequest{
		UserRoleId: model.UserRoleId,
	})

	if err != nil || data == nil {
		return nil, err
	}

	return &appsvc.CountUserResponse{
		TotalCount: data.TotalCount,
	}, err
}

func (s *appServer) SearchUser(ctx context.Context, model *appsvc.SearchUserRequest) (*appsvc.SearchUserResponse, error) {
	if model == nil {
		s.logger.Info("Search request is empty")
		return nil, errors.New("SEARCH_REQUEST_IS_EMPTY")
	}

	data, err := s.userApiClient.SearchUser(ctx, &user_sv.SearchUserRequest{
		SearchType: user_sv.SearchType(model.SearchType),
		SearchBy:   user_sv.UserAttribute(model.SearchBy),
		Keyword:    model.Keyword,
		UserRoleId: model.UserRoleId,
	})

	if err != nil || data == nil {
		return nil, err
	}

	var users []*appsvc.UserData
	for _, user := range data.Users {
		users = append(users, &appsvc.UserData{
			UserId:              user.UserId,
			UserName:            user.UserName,
			UserRole:            user.UserRole,
			UserGroups:          user.UserGroups,
			UserInfo:            user.UserInfo,
			LastModified:        user.LastModified,
			UserGroupIds:        user.UserGroupIds,
			RegisteredImageUrls: user.RegisteredImageUrls,
			ThumbnailImageUrl:   user.ThumbnailImageUrl,
			State:               appsvc.UserState(user.State),
			UserRoleId:          user.UserRoleId,
			IssuedDate:          user.IssuedDate,
			ExpiryDate:          user.ExpiryDate,
			ActivationDate:      user.ActivationDate,
			IsActive:            user.IsActive,
		})
	}

	return &appsvc.SearchUserResponse{
		Users: users,
		Total: data.Total,
	}, err
}

func (s *appServer) Update(ctx context.Context, model *appsvc.UpdateRequest) (*appsvc.UpdateResponse, error) {
	if model == nil {
		s.logger.Info("Update empty model")
		return nil, errors.New("Update empty model")
	}

	var images []*user_sv.Image

	if model.Images != nil {
		for _, image := range model.Images {
			if image != nil && image.Data != nil {
				images = append(images, &user_sv.Image{
					Data:    image.Data,
					ImageId: image.ImageId,
				})
			}
		}
	}

	data, err := s.userApiClient.Update(ctx, &user_sv.UpdateRequest{
		Images:         images,
		UserId:         model.UserId,
		UserName:       model.UserName,
		UserRoleId:     model.UserRoleId,
		UserGroupIds:   model.UserGroupIds,
		UserInfo:       model.UserInfo,
		ExpiryDate:     model.ExpiryDate,
		ActivationDate: model.ActivationDate,
		State:          user_sv.UserState(model.State),
	})

	if err != nil || data == nil {
		return nil, err
	}

	res := &appsvc.UpdateResponse{}
	copier.Copy(res, data)

	return res, err
}

func (s *appServer) Delete(ctx context.Context, model *appsvc.DeleteRequest) (*appsvc.DeleteResponse, error) {
	if model == nil {
		s.logger.Info("Delete empty model")
		return nil, errors.New("Delete empty model")
	}

	data, err := s.userApiClient.Delete(ctx,
		&user_sv.DeleteRequest{UserIds: model.UserIds})

	if err != nil || data == nil {
		return nil, err
	}

	return &appsvc.DeleteResponse{
		Code:    data.Code,
		Error:   data.Error,
		Message: data.Message,
		UserIds: data.UserIds,
	}, err
}
