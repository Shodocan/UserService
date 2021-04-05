package services

import (
	"github.com/Shodocan/UserService/internal/configs/engine"
	"github.com/Shodocan/UserService/internal/domain/entity"
)

//go:generate mockgen -destination user-repository_mock.go -package services . UserService
type UserService interface {
	Query(filters []entity.UserFilter, sortList []string, page, limit int) ([]entity.User, engine.Pagination, error)
	Create(data entity.User) (entity.User, error)
	Find(id string) (entity.User, error)
	Update(id string, user entity.User) (entity.User, error)
	PartialUpdate(id string, user entity.User) (entity.User, error)
	Delete(id string) error
}
