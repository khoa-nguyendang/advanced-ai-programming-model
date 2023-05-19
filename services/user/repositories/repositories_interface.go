package repositories

import (
	"aapi/pkg/logger"
	user_sv "aapi/services/user/protos/user/v1"
	"context"

	"aapi/shared/entities"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db     *sqlx.DB
	cache  *redis.Client
	logger logger.Logger
}

type UserRepository interface {
	// User
	EnrollUser(ctx context.Context, user *entities.User, image_URLs []string) (*entities.User, error)
	FindUserByUserId(ctx context.Context, company_code string, user_id string) (*entities.User, error)
	FindAllUser(ctx context.Context) []*entities.User
	UpdateUser(ctx context.Context, user *entities.User, image_URLs []string, company_code string) (*entities.User, error)
	DeleteUser(ctx context.Context, company_code string, user_id string) (bool, error)
	GetUser(ctx context.Context, model *user_sv.GetUserRequest, company_code string) (*user_sv.GetUserResponse, error)
	ValidateData(ctx context.Context, company_code string, group_ids []int64) (int32, string, string)
	GetEnrollmentImages(ctx context.Context, user_column_ids []int64) map[int64][]string
	CheckUserActivaton(ctx context.Context) (map[string][]string, map[string][]string)

	GetRoles(ctx context.Context) ([]*entities.Role, error)
	IsRoleValid(ctx context.Context, roleId int64) bool
	GetGroups(ctx context.Context) ([]*entities.Group, error)
	GetUserGroupsInfo(ctx context.Context, user_column_ids []int64) map[int64][]*entities.Group
}

// NewRepository func initializes a service
func NewRepository(db *sqlx.DB,
	logger logger.Logger,
	cache *redis.Client) UserRepository {
	return &userRepo{
		db:     db,
		logger: logger,
		cache:  cache,
	}
}
