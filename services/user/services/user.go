package services

import (
	"aapi/config"
	"aapi/pkg/logger"
	"aapi/pkg/minio_srv"
	mg "aapi/pkg/mongodb"
	ai_engine "aapi/services/user/protos/ai_engine/v1"
	faiss "aapi/services/user/protos/faiss/v1"
	user_sv "aapi/services/user/protos/user/v1"
	"aapi/services/user/repositories"
	"aapi/shared/constants"
	"aapi/shared/entities"
	httpclients "aapi/shared/http_clients"
	vm "aapi/shared/view_models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-co-op/gocron"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jinzhu/copier"
	"github.com/minio/minio-go/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
)

var company_code string = ""
var device_id string = ""

type service struct {
	logger         logger.Logger
	cfg            *config.Config
	tracer         opentracing.Tracer
	minioClient    *minio.Client
	kafkaProducer  *kafka.Producer
	engineClient   ai_engine.ImageAPIClient
	faissClient    faiss.VectorAPIClient
	userRepository repositories.UserRepository
	scheduler      *gocron.Scheduler

	user_sv.UnimplementedUserServer
}

type UserService interface {
	Enroll(ctx context.Context, model *user_sv.EnrollRequest) (*user_sv.EnrollResponse, error)
	Verify(ctx context.Context, model *user_sv.VerifyRequest) (*user_sv.VerifyResponse, error)
	GetUser(ctx context.Context, model *user_sv.GetUserRequest) (*user_sv.GetUserResponse, error)
	SearchUser(ctx context.Context, model *user_sv.SearchUserRequest) (*user_sv.SearchUserResponse, error)
	Update(ctx context.Context, model *user_sv.UpdateRequest) (*user_sv.UpdateResponse, error)
	Delete(ctx context.Context, model *user_sv.DeleteRequest) (*user_sv.DeleteResponse, error)

	Run(port int) error
}

const (
	USER_SERVICE_ENROLLMENT_BUCKET           = "user-service-enrollment"
	USER_SERVICE_ENROLLMENT_THUMNAILS_BUCKET = "user-service-enrollment-thumnails"
	USER_SERVICE_VERIFICATION_BUCKET         = "user-service-verification"

	USER_SERVICE_VERIFICATION_TOPIC = "user-service-verification-topic"
)

// NewService func initializes a service
func NewService(logger logger.Logger,
	cfg *config.Config,
	engineCnn *grpc.ClientConn,
	lightweightEngineCnn *grpc.ClientConn,
	faissCnn *grpc.ClientConn,
	userRepository repositories.UserRepository,
	trace opentracing.Tracer,
	minio *minio.Client,
	kafka_producer *kafka.Producer,
) UserService {
	return &service{
		logger:         logger,
		cfg:            cfg,
		engineClient:   ai_engine.NewImageAPIClient(engineCnn),
		faissClient:    faiss.NewVectorAPIClient(faissCnn),
		userRepository: userRepository,
		tracer:         trace,
		minioClient:    minio,
		kafkaProducer:  kafka_producer,
	}
}

// Run starts the server
func (s *service) Run(port int) error {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(s.tracer),
		),
		grpc.MaxRecvMsgSize(constants.MAX_MESSAGE_LENGTH),
		grpc.MaxSendMsgSize(constants.MAX_MESSAGE_LENGTH),
	)
	user_sv.RegisterUserServer(srv, s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.logger.Infof("failed to listen: %v", err)
	}

	minio_srv.CreateMinIOBucket(s.minioClient, USER_SERVICE_ENROLLMENT_BUCKET)
	minio_srv.CreateMinIOBucket(s.minioClient, USER_SERVICE_VERIFICATION_BUCKET)
	minio_srv.CreateMinIOBucket(s.minioClient, USER_SERVICE_ENROLLMENT_THUMNAILS_BUCKET)

	s.logger.Infof("Start Scheduler")
	s.scheduler = gocron.NewScheduler(time.Local)
	s.scheduler.StartAsync()

	return srv.Serve(lis)
}

func (s *service) Enroll(ctx context.Context, model *user_sv.EnrollRequest) (*user_sv.EnrollResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.Enroll")
	defer span.Finish()

	if model == nil {
		s.logger.Infof("Enroll empty model")
		return &user_sv.EnrollResponse{
			Code:    http.StatusBadRequest,
			Message: "Enroll with empty request",
			Error:   "ENROLL_WITH_EMPTY_REQUEST",
		}, nil
	}

	s.logger.Infof("Enroll data: user_id: %v, images count: %v ", model.UserId, len(model.Images))

	existing_user, err := s.userRepository.FindUserByUserId(ctx, company_code, model.UserId)
	s.logger.Infof("Existing search %v", existing_user)
	if err != nil {
		return nil, err
	}

	if existing_user != nil {
		return &user_sv.EnrollResponse{
			Code:    http.StatusBadRequest,
			Message: "User Id is already registered",
			Error:   "USER_ID_DUPLICATED",
		}, nil
	}

	if model.ExpiryDate == 0 {
		model.ExpiryDate = 8640000000
	}

	var image_URLs []string
	// Enroll by images
	if model.Images != nil && len(model.Images) > 0 {
		images := make([]*ai_engine.Image, len(model.Images))
		for idx, image := range model.Images {
			if image == nil {
				s.logger.Infof("image is null")
				continue
			}
			s.logger.Infof("image size: %v", len(image.Data))
			imageId := fmt.Sprintf("%s_%d", model.UserId, idx)
			image.ImageId = imageId

			images[idx] = &ai_engine.Image{
				Data:    image.Data,
				ImageId: image.ImageId,
			}
		}
		engineResponse, err := s.engineClient.Enroll(ctx,
			&ai_engine.ImageEnrollmentRequest{
				UserUid: model.UserId,
				Images:  images,
			},
		)

		if err != nil {
			s.logger.Errorf("Extract feature error: %v", err)
			return &user_sv.EnrollResponse{
				Code:    http.StatusInternalServerError,
				Message: "Extract feature error",
				Error:   "INTERNAL_SERVER_ERROR",
			}, nil
		}

		s.logger.Infof("Got result from Engine: %v", engineResponse)
		if engineResponse == nil || err != nil {
			s.logger.Errorf("Got error from engine: %v", err.Error())
			return nil, err
		}

		if engineResponse.Code != ai_engine.StatusCode_SUCCESS {
			switch engineResponse.Code {
			case ai_engine.StatusCode_BAD_IMAGE:
				return &user_sv.EnrollResponse{
					Code:    http.StatusBadRequest,
					Message: engineResponse.Message,
					Error:   "BAD_IMAGE",
				}, nil

			case ai_engine.StatusCode_ENGINE_ERROR:
				return &user_sv.EnrollResponse{
					Code:    http.StatusInternalServerError,
					Message: "AI Engine error",
					Error:   "INTERNAL_SERVER_ERROR",
				}, nil

			case ai_engine.StatusCode_IMAGE_REGISTERED_FOR_ANOTHER_USER:
				return &user_sv.EnrollResponse{
					Code:    http.StatusBadRequest,
					Message: "The image input is already registered for another user Id",
					Error:   "IMAGE_REGISTERED_FOR_ANOTHER_USER",
				}, nil

			default:
				return &user_sv.EnrollResponse{
					Code:    http.StatusInternalServerError,
					Message: "AI Engine error",
					Error:   "INTERNAL_SERVER_ERROR",
				}, nil
			}
		}

		currentTime := time.Now().Local().Format(time.RFC3339Nano)
		currentDate := currentTime[0:10]
		for _, image := range model.Images {
			target_file_name := currentDate + "/" + model.UserId + "/" + image.ImageId + ".jpg"
			file_reader := bytes.NewReader(image.Data)
			minio_srv.UploadFile(s.minioClient, USER_SERVICE_ENROLLMENT_BUCKET, target_file_name, file_reader, file_reader.Size())
			image_URLs = append(image_URLs, target_file_name)
		}

		user, err := s.saveUserData(ctx, model, company_code, image_URLs, image_URLs[0], constants.UserState(constants.USER_STATE_ENABLED))
		s.logger.Infof("Enroll wroten data: %#v - err: %v", user, err)

		return &user_sv.EnrollResponse{
			Code:    http.StatusOK,
			Message: "Enroll user successfully",
			Error:   "",
			State:   user_sv.UserState(user.UserState),
		}, nil
	}

	return &user_sv.EnrollResponse{
		Code:    http.StatusBadRequest,
		Message: "Request has no image",
		Error:   "REQUEST_HAS_NO_IMAGE",
	}, nil
}

func (s *service) Verify(ctx context.Context, model *user_sv.VerifyRequest) (*user_sv.VerifyResponse, error) {
	// Verify by images
	if model.Images != nil && len(model.Images) > 0 {
		s.logger.Infof("UserService.Verify images")

		images := make([]*ai_engine.Image, len(model.Images))
		currentTime := time.Now().Local().Format(time.RFC3339Nano)
		currentDate := currentTime[0:10]
		s.logger.Infof(currentTime)

		for idx, image := range model.Images {
			if image == nil {
				s.logger.Infof("image is null")
				continue
			}

			if constants.ENABLED == os.Getenv(constants.VERIFICATION_IMAGE_LOG) {
				image.ImageId = fmt.Sprintf("%s/Device_%s/Session_%s/Image_%d_%s", currentDate, device_id, image.ImageId, idx, currentTime)
			}

			images[idx] = &ai_engine.Image{
				Data:    image.Data,
				ImageId: image.ImageId,
			}
		}
		searchResponse, err := s.engineClient.Search(ctx,
			&ai_engine.ImageSearchRequest{
				Images: images,
			},
		)
		s.logger.Infof("Got result from engineClient: %v", searchResponse)

		if err != nil {
			s.logger.Errorf("got err from ai_engine.Search: %v", err)
			// TODO: disabled due to handle at Client side
			// s.handleVerificationFailed(ctx, device_id, company_code, model.Images[0])
			return &user_sv.VerifyResponse{
				Code:    http.StatusInternalServerError,
				Message: "Extract feature error",
				Error:   "INTERNAL_SERVER_ERROR",
			}, nil
		}

		if !(searchResponse.Code == ai_engine.StatusCode_SUCCESS || searchResponse.Code == ai_engine.StatusCode_FOUND) {
			// TODO: disabled due to handle at Client side
			//s.handleVerificationFailed(ctx, device_id, company_code, model.Images[0])
			switch searchResponse.Code {
			case ai_engine.StatusCode_BAD_IMAGE:
				return &user_sv.VerifyResponse{
					Code:    http.StatusBadRequest,
					Message: searchResponse.Message,
					Error:   "BAD_IMAGE",
				}, nil

			case ai_engine.StatusCode_ENGINE_ERROR:
				return &user_sv.VerifyResponse{
					Code:    http.StatusInternalServerError,
					Message: "AI Engine error",
					Error:   "INTERNAL_SERVER_ERROR",
				}, nil

			case ai_engine.StatusCode_NOT_FOUND:
				return &user_sv.VerifyResponse{
					Code:    http.StatusBadRequest,
					Message: "Cannot found any user id for the given images",
					Error:   "USER_NOT_FOUND",
				}, nil

			default:
				return &user_sv.VerifyResponse{
					Code:    http.StatusInternalServerError,
					Message: "AI Engine error",
					Error:   "INTERNAL_SERVER_ERROR",
				}, nil
			}
		}

		// Get user from db
		user, err := s.userRepository.FindUserByUserId(ctx, company_code, searchResponse.UserUid)

		if err != nil || user == nil {
			s.logger.Errorf("got err from DB or user not found: %v", err)
			// TODO: disabled due to handle at Client side
			//s.handleVerificationFailed(ctx, device_id, company_code, model.Images[0])
			return &user_sv.VerifyResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot found any user id for the given images",
				Error:   "USER_NOT_FOUND",
			}, nil
		}

		if err != nil {
			s.logger.Errorf("Verify.err: %v", err)
			return &user_sv.VerifyResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot found any user id for the given images",
				Error:   "USER_NOT_FOUND",
			}, nil
		}

		is_active := false
		current_time := time.Now().UTC().Unix()
		if (user.UserRoleId == 4 && user.UserState == constants.USER_STATE_ENABLED) || (user.UserRoleId == 5 && user.UserState == constants.USER_STATE_ENABLED && user.ActivationDate < current_time && current_time < user.ExpiryDate) {
			is_active = true
		}
		// TODO: disabled due to handle at Client side
		//s.handleVerificationSuccess(ctx, user.UserId, user.UserName, device_id, company_code, is_active, model.Images[0])

		roles, _ := s.userRepository.GetRoles(ctx)
		var role_name string

		for _, role := range roles {
			if role.Id == user.UserRoleId {
				role_name = role.Name
				break
			}
		}

		// user_groups := s.userRepository.GetUserGroupsInfo(ctx, []int64{user.Id})

		group_names := []string{}
		// user_group_info := user_groups[user.Id]
		// for _, group_info := range user_group_info {
		// 	group_names = append(group_names, group_info.Name)
		// }

		return &user_sv.VerifyResponse{
			Code:              http.StatusOK,
			Message:           "User identification successfully",
			UserId:            user.UserId,
			UserName:          user.UserName,
			UserRole:          role_name,
			UserGroups:        group_names,
			IssuedDate:        user.IssuedDate,
			ActivationDate:    user.ActivationDate,
			ExpiryDate:        user.ExpiryDate,
			State:             user_sv.UserState(user.UserState),
			ThumbnailImageUrl: minio_srv.GetTmpURL(s.minioClient, USER_SERVICE_ENROLLMENT_BUCKET, user.ThumbnailImageUrl),
			UserInfo:          user.UserInfo,
			IsActive:          is_active,
			ImageId:           searchResponse.ImageId,
			Score:             searchResponse.Score,
			LastModified:      user.LastModified,
		}, err
	}
	// TODO: disabled due to handle at Client side
	//s.handleVerificationFailed(ctx, device_id, company_code, model.Images[0])
	return &user_sv.VerifyResponse{
		Code:    http.StatusBadRequest,
		Message: "Request has no image",
		Error:   "REQUEST_HAS_NO_IMAGE",
	}, nil
}

func (s *service) GetUser(ctx context.Context, model *user_sv.GetUserRequest) (*user_sv.GetUserResponse, error) {
	s.logger.Infof("GetUser triggered: %#v", model)

	if model == nil {
		s.logger.Info("GetUser with empty request")
		return nil, errors.New("GetUser with empty request")
	}

	res, err := s.userRepository.GetUser(ctx, model, company_code)
	if res != nil && len(res.Users) > 0 {
		for i := range res.Users {
			for j := range res.Users[i].RegisteredImageUrls {
				res.Users[i].RegisteredImageUrls[j] = minio_srv.GetTmpURL(s.minioClient, USER_SERVICE_ENROLLMENT_BUCKET, res.Users[i].RegisteredImageUrls[j])
			}
			res.Users[i].ThumbnailImageUrl = minio_srv.GetTmpURL(s.minioClient, USER_SERVICE_ENROLLMENT_BUCKET, res.Users[i].ThumbnailImageUrl)
		}
	}
	return res, err
}

func (s *service) CountUser(ctx context.Context, model *user_sv.CountUserRequest) (*user_sv.CountUserResponse, error) {
	if model == nil {
		s.logger.Info("CountUser request is empty")
		return nil, errors.New("COUNT_REQUEST_IS_EMPTY")
	}

	mongoConfig := vm.MongoDbConfig{
		UserName:   s.cfg.MongoAtlas.MongodbUser,
		Password:   s.cfg.MongoAtlas.MongodbPassword,
		Url:        s.cfg.MongoAtlas.MongodbHost,
		Database:   s.cfg.MongoAtlas.MongodbDbname,
		Collection: "users",
	}

	roles, _ := s.userRepository.GetRoles(ctx)
	var role_name string

	for _, role := range roles {
		if role.Id == model.UserRoleId {
			role_name = role.Name
			break
		}
	}

	filter := bson.M{"company": company_code, "role": role_name}
	if role_name == "" {
		filter = bson.M{"company": company_code}
	}

	total_count, err := mg.CloudCollectionCount(mongoConfig, filter)

	s.logger.Infof("CountUser : total_count %v err %v", total_count, err)

	return &user_sv.CountUserResponse{
		TotalCount: total_count,
	}, err
}

func (s *service) SearchUser(ctx context.Context, model *user_sv.SearchUserRequest) (*user_sv.SearchUserResponse, error) {
	return nil, errors.Errorf("Unimplemented yet")
}

func (s *service) Update(ctx context.Context, model *user_sv.UpdateRequest) (*user_sv.UpdateResponse, error) {
	s.logger.Infof("Update user %v", model.UserId)
	if model == nil {
		s.logger.Info("Update user with empty model")
		return &user_sv.UpdateResponse{
			Code:    http.StatusBadRequest,
			Message: "Request body is missing",
			Error:   "Empty request body",
		}, nil
	}

	if !s.userRepository.IsRoleValid(ctx, model.UserRoleId) {
		return &user_sv.UpdateResponse{
			Code:    http.StatusBadRequest,
			Message: "Role is not valid",
			Error:   "ROLE_ID_NOT_VALID",
		}, nil
	}

	existing_user, err := s.userRepository.FindUserByUserId(ctx, company_code, model.UserId)

	if existing_user == nil || err != nil {
		s.logger.Info("Update user to non-existing model")
		return &user_sv.UpdateResponse{
			Code:    http.StatusBadRequest,
			Message: "User Id doesn't exsit",
			Error:   "USER_ID_DOES_NOT_EXIST",
		}, nil
	}

	current_time := time.Now().UTC().Unix()

	if existing_user.UserRoleId == 5 && model.UserRoleId == 4 {
		model.ActivationDate = 0
		model.ExpiryDate = 8640000000
	}

	is_active := false
	if (model.UserRoleId == 4 && model.State == constants.USER_STATE_ENABLED) || (model.UserRoleId == 5 && model.State == constants.USER_STATE_ENABLED && model.ActivationDate < current_time && current_time < model.ExpiryDate) {
		s.logger.Infof("ActivationDate %v - ExpiryDate %v - current time %v", model.ActivationDate, model.ExpiryDate, current_time)
		is_active = true
	}

	// Update images
	var new_image_URLs []string
	registered_image_URLs := s.userRepository.GetEnrollmentImages(ctx, []int64{existing_user.Id})[existing_user.Id]
	current_length := len(registered_image_URLs)
	if model.Images != nil && len(model.Images) > 0 {
		s.logger.Infof("Update images - Add : %#v more", len(model.GetImages()))

		images := make([]*ai_engine.Image, len(model.Images))
		for idx, image := range model.Images {
			if image == nil {
				s.logger.Infof("image is null")
				continue
			}
			s.logger.Infof("image size: %v", len(image.Data))
			imageId := fmt.Sprintf("%s_%d", model.UserId, idx+current_length)
			image.ImageId = imageId

			images[idx] = &ai_engine.Image{
				Data:    image.Data,
				ImageId: image.ImageId,
			}
		}

		// Update on FAISS
		engineResponse, err := s.engineClient.Enroll(ctx,
			&ai_engine.ImageEnrollmentRequest{
				UserUid: model.UserId,
				Images:  images,
			},
		)

		if err != nil {
			s.logger.Info("Extract feature error")
			return &user_sv.UpdateResponse{
				Code:    http.StatusInternalServerError,
				Message: "Extract feature error",
				Error:   "INTERNAL_SERVER_ERROR",
			}, nil
		}

		s.logger.Infof("Got result from Engine: %v", engineResponse)
		if engineResponse == nil || err != nil {
			s.logger.Errorf("Got error from engine: %v", err.Error())
			return nil, err
		}

		if engineResponse.Code != ai_engine.StatusCode_SUCCESS {
			switch engineResponse.Code {
			case ai_engine.StatusCode_BAD_IMAGE:
				return &user_sv.UpdateResponse{
					Code:    http.StatusBadRequest,
					Message: engineResponse.Message,
					Error:   "BAD_IMAGE",
				}, nil

			case ai_engine.StatusCode_ENGINE_ERROR:
				return &user_sv.UpdateResponse{
					Code:    http.StatusInternalServerError,
					Message: "AI Engine error",
					Error:   "INTERNAL_SERVER_ERROR",
				}, nil

			case ai_engine.StatusCode_IMAGE_REGISTERED_FOR_ANOTHER_USER:
				return &user_sv.UpdateResponse{
					Code:    http.StatusBadRequest,
					Message: "The image input is already registered for another user Id",
					Error:   "IMAGE_REGISTERED_FOR_ANOTHER_USER",
				}, nil

			default:
				return &user_sv.UpdateResponse{
					Code:    http.StatusInternalServerError,
					Message: "AI Engine error",
					Error:   "INTERNAL_SERVER_ERROR",
				}, nil
			}
		}

		// Update to light weight engine with registered images + new images
		for idx, registered_image_URL := range registered_image_URLs {
			file_bytes := minio_srv.DownloadFile(s.minioClient, USER_SERVICE_ENROLLMENT_BUCKET, registered_image_URL)
			if file_bytes != nil {
				registered_image := ai_engine.Image{
					Data:    file_bytes,
					ImageId: fmt.Sprintf("%s_%d", model.UserId, idx),
				}
				images = append(images, &registered_image)
			} else {
				return &user_sv.UpdateResponse{
					Code:    http.StatusInternalServerError,
					Message: "Cannot get registered images",
					Error:   "INTERNAL_SERVER_ERROR",
				}, nil
			}
		}

		// Delete registered data
		delAtlasResponse, atlas_err := s.faissClient.Delete(ctx,
			&faiss.VectorDeletionRequest{
				UserUid: model.UserId,
			},
		)
		s.logger.Info(delAtlasResponse)
		if atlas_err != nil {
			s.logger.Errorf("DeleteUser err: %v", atlas_err)

			// Rollback
			delResponse, faiss_err := s.faissClient.Delete(ctx,
				&faiss.VectorDeletionRequest{
					UserUid: model.UserId,
				},
			)
			s.logger.Infof("Res %v Err %v", delResponse, faiss_err)
			return &user_sv.UpdateResponse{
				Code:    http.StatusInternalServerError,
				Message: "Cannot remove registered image on Mongo Atlas",
				Error:   "INTERNAL_SERVER_ERROR",
			}, nil
		}

		lightweightEngineResponse, _ := s.engineClient.Enroll(ctx,
			&ai_engine.ImageEnrollmentRequest{
				Images: images,
			},
		)
		s.logger.Infof("Got result from lightweight engine: %v", lightweightEngineResponse)
		// TO-DO: Should let user service update MongoDB to control logic here easily
		if lightweightEngineResponse == nil || lightweightEngineResponse.Code != ai_engine.StatusCode_SUCCESS {
			// Request Robus AI Engine service remove feature in DB
			s.logger.Infof("Images of %s failed to register with lightweght model", model.UserId)
			res, err := s.faissClient.Delete(ctx, &faiss.VectorDeletionRequest{UserUid: model.UserId})
			s.logger.Infof("Delete on FAISS result: %v err: %v", res, err)
			if lightweightEngineResponse == nil {
				return &user_sv.UpdateResponse{
					Code:    http.StatusInternalServerError,
					Message: "AI Engine error",
					Error:   "INTERNAL_SERVER_ERROR",
				}, nil
			}
			return &user_sv.UpdateResponse{
				Code:    http.StatusBadRequest,
				Message: lightweightEngineResponse.Message,
				Error:   "BAD_IMAGE",
			}, nil
		}

		currentTime := time.Now().Local().Format(time.RFC3339Nano)
		currentDate := currentTime[0:10]

		// Upload to minio the new images
		for _, image := range model.Images {
			target_file_name := currentDate + "/" + model.UserId + "/" + image.ImageId + ".jpg"
			file_reader := bytes.NewReader(image.Data)
			minio_srv.UploadFile(s.minioClient, USER_SERVICE_ENROLLMENT_BUCKET, target_file_name, file_reader, file_reader.Size())
			s.logger.Infof("Uploaded new image %v", target_file_name)
			new_image_URLs = append(new_image_URLs, target_file_name)

		}
	} else {
		roles, _ := s.userRepository.GetRoles(ctx)
		var role_name string

		for _, role := range roles {
			if role.Id == model.UserRoleId {
				role_name = role.Name
				break
			}
		}
		s.updateUserOnCloud(company_code, model.UserId, model.UserName, role_name, model.UserInfo, is_active)
	}

	entt := &entities.User{}
	copier.Copy(entt, &model)
	entt.Id = existing_user.Id
	entt.UserState = constants.UserState(model.State)
	entt.CompanyCode = existing_user.CompanyCode
	entt.ThumbnailImageUrl = existing_user.ThumbnailImageUrl
	user, err := s.userRepository.UpdateUser(ctx, entt, new_image_URLs, company_code)

	if err != nil {
		s.logger.Infof("Update user error: %v", err)
		return &user_sv.UpdateResponse{
			Code:    http.StatusInternalServerError,
			Message: "Cannot update to database",
			Error:   "INTERNAL_SERVER_ERROR",
		}, nil
	}

	s.logger.Infof("Updated user: %#v", user)

	return &user_sv.UpdateResponse{
		Code:    http.StatusOK,
		Message: "User updated successfully",
	}, nil
}

// Delete only delete 1 item
func (s *service) Delete(ctx context.Context, model *user_sv.DeleteRequest) (*user_sv.DeleteResponse, error) {
	s.logger.Infof("Delete user : %#v", model)
	if model == nil {
		s.logger.Info("Delete user with empty request")
		return nil, errors.New("Delete user with empty request")
	}

	var deleted_user_ids []string

	for _, user_id := range model.UserIds {
		existing_user, find_err := s.userRepository.FindUserByUserId(ctx, company_code, user_id)

		if existing_user == nil || find_err != nil {
			s.logger.Info("Delete a non-existing user")
			continue
		}

		delResponse, faiss_err := s.faissClient.Delete(ctx,
			&faiss.VectorDeletionRequest{
				UserUid: user_id,
			},
		)
		s.logger.Info(delResponse)

		if faiss_err == nil && delResponse.Code == faiss.StatusCode_SUCCESS {
			deleted, del_err := s.userRepository.DeleteUser(ctx, company_code, existing_user.UserId)

			if del_err == nil && deleted {
				deleted_user_ids = append(deleted_user_ids, user_id)
			} else {
				s.logger.Errorf("DeleteUser err: %v", del_err)
			}
		}
	}

	if len(deleted_user_ids) == 0 {
		return &user_sv.DeleteResponse{
			Code:    http.StatusBadRequest,
			Message: "Cannot delete by the given user ids",
			Error:   "CANNOT_DELETE_USER_IDS",
			UserIds: deleted_user_ids,
		}, nil
	}
	return &user_sv.DeleteResponse{
		Code:    http.StatusOK,
		Message: "User deleted successfully",
		UserIds: deleted_user_ids,
	}, nil
}

func (s *service) saveUserData(ctx context.Context, model *user_sv.EnrollRequest, company_code string, image_URLs []string, thumnail_image_url string, user_state constants.UserState) (*entities.User, error) {
	ett := &entities.User{}
	copier.Copy(ett, model)
	ett.UserState = user_state
	ett.CompanyCode = company_code
	ett.ThumbnailImageUrl = thumnail_image_url
	// Write user info to database
	user, err := s.userRepository.EnrollUser(ctx, ett, image_URLs)

	return user, err
}

// handleVerificationFailed
func (s *service) handleVerificationFailed(ctx context.Context, deviceUUid, company_code string, img *user_sv.Image) {
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
		return
	}

	s.handleVerification(ctx, string(jsondata), deviceUUid, company_code)
}

func (s *service) handleVerificationSuccess(ctx context.Context,
	userId,
	userName,
	deviceUUid,
	company_code string,
	userState bool,
	img *user_sv.Image) {

	imgPath, preSignedUrl := s.uploadImageToMinio(img, company_code, true)

	dt := &vm.MatchingInput{
		UserId:       userId,
		UserName:     userName,
		Temperature:  "",
		DeviceName:   deviceUUid,
		ImagePathS3:  imgPath,
		PreSignedUrl: preSignedUrl,
		VerifyStatus: true,
		UserStatus:   userState,
	}
	matchingJson, err := json.Marshal(dt)
	if err != nil {
		s.logger.Errorf("handleVerificationSuccess got error when marshal MatchingInput struct: %v", err)
		return
	}
	s.handleVerification(ctx, string(matchingJson), deviceUUid, company_code)
}

func (s *service) handleVerification(ctx context.Context, userJson, deviceUUid, company_code string) {
	message := vm.CommonMessage{
		MessageType:  constants.VERIFICATION_EVENT,
		From:         "Device" + deviceUUid,
		To:           "WebAdmin",
		Topic:        company_code,
		CompanyCode:  company_code,
		Body:         userJson,
		CreatedDate:  time.Now().UTC().UnixMilli(),
		LastModified: time.Now().UTC().UnixMilli(),
		Acknowledge:  false,
	}
	s.broadcastVerification(message)
}

func (s *service) broadcastVerification(message vm.CommonMessage) {
	s.logger.Infof("broadcastVerification: %v", message)
	url := s.cfg.PublisherServer + "/send-message"
	responseMessage, err := httpclients.SendHttpPost2(url, constants.Annonymous, message, "", "", "")
	if err != nil {
		s.logger.Errorf("broadcastVerification got error: %v", err)
	}

	s.logger.Infof("broadcastVerification.success: %v", string(responseMessage))
}

// uploadImageToMinio upload a picture to minio then return path, presigned url
func (s *service) uploadImageToMinio(img *user_sv.Image, company_code string, status bool) (path string, presignedUrl string) {
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

func (s *service) doDailyTask() {
	s.logger.Infof("doDailyTask: BEGIN")
	s.logger.Infof("doDailyTask: Update UserState")
	activated_users, deactivated_users := s.userRepository.CheckUserActivaton(context.Background())
	for company_code, user_ids := range activated_users {
		for _, user_id := range user_ids {
			s.updateUserStateOnCloud(company_code, user_id, true)
		}
	}
	for company_code, user_ids := range deactivated_users {
		for _, user_id := range user_ids {
			s.updateUserStateOnCloud(company_code, user_id, false)
		}
	}
	s.logger.Infof("doDailyTask: FINISHED")
}

func (s *service) updateUserStateOnCloud(company_code string, user_id string, is_active bool) {
	mongoConfig := vm.MongoDbConfig{
		UserName:   s.cfg.MongoAtlas.MongodbUser,
		Password:   s.cfg.MongoAtlas.MongodbPassword,
		Url:        s.cfg.MongoAtlas.MongodbHost,
		Database:   s.cfg.MongoAtlas.MongodbDbname,
		Collection: "users",
	}

	filter_user := entities.MongoUserBase{
		EmployeeId: user_id,
		Company:    company_code,
	}

	update := bson.M{"$set": bson.M{"isActive": is_active}}

	update_number, mongo_err := mg.CloudUpdateOne(mongoConfig, filter_user, update)
	s.logger.Infof("Sync state for user: %v Company: %v - DB Log: Update number %v - Err: %v", user_id, company_code, update_number, mongo_err)
}

func (s *service) updateUserOnCloud(company_code string, user_id string, user_name string, role string, user_info string, is_active bool) {
	mongoConfig := vm.MongoDbConfig{
		UserName:   s.cfg.MongoAtlas.MongodbUser,
		Password:   s.cfg.MongoAtlas.MongodbPassword,
		Url:        s.cfg.MongoAtlas.MongodbHost,
		Database:   s.cfg.MongoAtlas.MongodbDbname,
		Collection: "users",
	}

	filter_user := entities.MongoUserBase{
		EmployeeId: user_id,
		Company:    company_code,
	}

	update := bson.M{"$set": bson.M{"isActive": is_active, "name": user_name, "role": role, "userInfo": user_info}}

	update_number, mongo_err := mg.CloudUpdateOne(mongoConfig, filter_user, update)
	s.logger.Infof("Sync state for user: %v Company: %v - DB Log: Update number %v - Err: %v", user_id, company_code, update_number, mongo_err)
}
