package usecase

import (
	"log"

	"github.com/Shodocan/UserService/internal/configs"
	"github.com/Shodocan/UserService/internal/configs/engine"
	"github.com/Shodocan/UserService/internal/domain/entity"
	"github.com/Shodocan/UserService/internal/helpers"
	"github.com/Shodocan/UserService/internal/services"
)

type UserCase struct {
	config  *configs.EnvVarConfig
	service services.UserService
	log     *log.Logger
}

func NewUserCase(config *configs.EnvVarConfig,
	service services.UserService,
	log *log.Logger,
) *UserCase {
	return &UserCase{config: config, service: service, log: log}
}

func (cs UserCase) Search(filters []entity.UserFilter, sort []string, limit int, page int) ([]entity.User, engine.Pagination, error) {
	usrs, paginate, err := cs.service.Query(filters, sort, page, limit)
	if err != nil {
		cs.log.Println(err)
		return nil, engine.Pagination{Pages: 0, Total: 0, PageSize: 0}, engine.ErrInternalFailure()
	}

	for i := range usrs {
		usrs[i].Password = ""
	}

	return usrs, paginate, nil
}

func (cs UserCase) ValidatePassword(id, password string) error {
	usr, err := cs.service.Find(id)
	if err != nil {
		cs.log.Println(err)
		return engine.ErrInternalFailure()
	}

	compareErr := helpers.ComparePasswordToHash(usr.Password, password)
	if compareErr != nil {
		return compareErr
	}

	return nil
}

func (cs UserCase) Find(id string) (entity.User, error) {
	usr, err := cs.service.Find(id)
	if err != nil {
		cs.log.Println(err)
		return entity.User{}, engine.ErrInternalFailure()
	}

	usr.Password = "" // must never return password to the user
	return usr, nil
}

func (cs UserCase) Create(user entity.User) (entity.User, error) {
	passHash, hasErr := helpers.HashPassword(user.Password)
	if hasErr != nil {
		return entity.User{}, hasErr
	}

	user.Password = passHash

	usr, err := cs.service.Create(user)
	if err != nil {
		cs.log.Println(err)
		return entity.User{}, engine.ErrInternalFailure()
	}

	usr.Password = "" // must never return password to the user
	return usr, nil
}

func (cs UserCase) Update(id string, user entity.User) (entity.User, error) {
	if user.Password != "" {
		passHash, hasErr := helpers.HashPassword(user.Password)
		if hasErr != nil {
			return entity.User{}, hasErr
		}

		user.Password = passHash
	}
	usr, err := cs.service.Update(id, user)
	if err != nil {
		cs.log.Println(err)
		return entity.User{}, engine.ErrInternalFailure()
	}

	usr.Password = "" // must never return password to the user
	return usr, nil
}

func (cs UserCase) PartialUpdate(id string, user entity.User) (entity.User, error) {
	if user.Password != "" {
		passHash, hasErr := helpers.HashPassword(user.Password)
		if hasErr != nil {
			return entity.User{}, hasErr
		}

		user.Password = passHash
	}
	usr, err := cs.service.PartialUpdate(id, user)
	if err != nil {
		cs.log.Println(err)
		return entity.User{}, engine.ErrInternalFailure()
	}

	usr.Password = "" // must never return password to the user
	return usr, nil
}

func (cs UserCase) Delete(id string) error {
	err := cs.service.Delete(id)
	if err != nil {
		cs.log.Println(err)
		return engine.ErrInternalFailure()
	}
	return nil
}
