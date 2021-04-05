package usecase

import (
	"fmt"
	"testing"

	"github.com/Shodocan/UserService/internal/configs"
	"github.com/Shodocan/UserService/internal/configs/engine"
	"github.com/Shodocan/UserService/internal/domain/entity"
	"github.com/Shodocan/UserService/internal/helpers"
	"github.com/Shodocan/UserService/internal/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func getConfig(t *testing.T) *configs.EnvVarConfig {
	config, err := configs.GetEnvConfig()
	if err != nil {
		t.Log(err)
		return config
	}
	return config
}
func TestSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	users := []entity.User{
		{ID: "123", Name: "Walisson Casonatto", Email: "wdcasonatto@gmail.com", Password: "32131"},
		{ID: "456", Name: "Samara Casonatto", Email: "wdcasonatto@gmail.com", Password: "509872"},
	}

	pagination := engine.NewPagination(2, 2, 10)

	userServiceMock := services.NewMockUserService(ctrl)
	userServiceMock.EXPECT().Query(gomock.Any(), []string{}, 10, 1).Return(users, pagination, nil)

	retreivedUsers, page, err := NewUserCase(getConfig(t), userServiceMock, configs.NewLog()).Search([]entity.UserFilter{}, []string{}, 1, 10)
	assert.Nil(t, err)
	assert.Equal(t, page, pagination)
	assert.Equal(t, retreivedUsers, users)
	for _, user := range retreivedUsers {
		assert.Empty(t, user.Password)
	}
}

func TestSearchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userServiceMock := services.NewMockUserService(ctrl)
	userServiceMock.EXPECT().Query(gomock.Any(), []string{}, 10, 1).Return(nil, engine.Pagination{}, fmt.Errorf("adfasdf"))

	retreivedUsers, page, err := NewUserCase(getConfig(t), userServiceMock, configs.NewLog()).Search([]entity.UserFilter{}, []string{}, 1, 10)
	assert.NotNil(t, err)
	assert.Equal(t, page.Total, 0)
	assert.Len(t, retreivedUsers, 0)
}

func TestFind(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := entity.User{ID: "123", Name: "Walisson Casonatto", Email: "wdcasonatto@gmail.com", Password: "32131"}

	userServiceMock := services.NewMockUserService(ctrl)
	userServiceMock.EXPECT().Find("123").Return(user, nil)

	retreivedUser, err := NewUserCase(getConfig(t), userServiceMock, configs.NewLog()).Find("123")
	assert.Nil(t, err)
	assert.Equal(t, retreivedUser.ID, "123")
	assert.Equal(t, retreivedUser.Name, "Walisson Casonatto")
	assert.Empty(t, retreivedUser.Password)
}

func TestFindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := entity.User{}

	userServiceMock := services.NewMockUserService(ctrl)
	userServiceMock.EXPECT().Find("123").Return(user, fmt.Errorf("123"))

	_, err := NewUserCase(getConfig(t), userServiceMock, configs.NewLog()).Find("123")
	assert.NotNil(t, err)
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := entity.User{ID: "123", Name: "Walisson Casonatto", Email: "wdcasonatto@gmail.com", Password: "32131"}

	userServiceMock := services.NewMockUserService(ctrl)
	userServiceMock.EXPECT().Create(gomock.Any()).Return(user, nil)

	retreivedUser, err := NewUserCase(getConfig(t), userServiceMock, configs.NewLog()).Create(user)
	assert.Nil(t, err)
	assert.Equal(t, retreivedUser.ID, "123")
	assert.Equal(t, retreivedUser.Name, "Walisson Casonatto")
	assert.Empty(t, retreivedUser.Password)
}

func TestCreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := entity.User{ID: "123", Name: "Walisson Casonatto", Email: "wdcasonatto@gmail.com", Password: "32131"}

	userServiceMock := services.NewMockUserService(ctrl)
	userServiceMock.EXPECT().Create(gomock.Any()).Return(entity.User{}, fmt.Errorf("123"))

	_, err := NewUserCase(getConfig(t), userServiceMock, configs.NewLog()).Create(user)
	assert.NotNil(t, err)
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := entity.User{ID: "123", Name: "Walisson Casonatto", Email: "wdcasonatto@gmail.com", Password: "32131"}

	userServiceMock := services.NewMockUserService(ctrl)
	userServiceMock.EXPECT().Update("123", gomock.Any()).Return(user, nil)

	retreivedUser, err := NewUserCase(getConfig(t), userServiceMock, configs.NewLog()).Update("123", user)
	assert.Nil(t, err)
	assert.Equal(t, retreivedUser.ID, "123")
	assert.Equal(t, retreivedUser.Name, "Walisson Casonatto")
	assert.Empty(t, retreivedUser.Password)
}

func TestUpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := entity.User{ID: "123", Name: "Walisson Casonatto", Email: "wdcasonatto@gmail.com", Password: "32131"}

	userServiceMock := services.NewMockUserService(ctrl)
	userServiceMock.EXPECT().Update("123", gomock.Any()).Return(entity.User{}, fmt.Errorf("123"))

	_, err := NewUserCase(getConfig(t), userServiceMock, configs.NewLog()).Update("123", user)
	assert.NotNil(t, err)
}

func TestPartialUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := entity.User{ID: "123", Name: "Walisson Casonatto", Email: "wdcasonatto@gmail.com", Password: "32131"}

	userServiceMock := services.NewMockUserService(ctrl)
	userServiceMock.EXPECT().PartialUpdate("123", gomock.Any()).Return(user, nil)

	retreivedUser, err := NewUserCase(getConfig(t), userServiceMock, configs.NewLog()).PartialUpdate("123", user)
	assert.Nil(t, err)
	assert.Equal(t, retreivedUser.ID, "123")
	assert.Equal(t, retreivedUser.Name, "Walisson Casonatto")
	assert.Empty(t, retreivedUser.Password)
}

func TestPartialUpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := entity.User{ID: "123", Name: "Walisson Casonatto", Email: "wdcasonatto@gmail.com", Password: "32131"}

	userServiceMock := services.NewMockUserService(ctrl)
	userServiceMock.EXPECT().PartialUpdate("123", gomock.Any()).Return(entity.User{}, fmt.Errorf("123"))

	_, err := NewUserCase(getConfig(t), userServiceMock, configs.NewLog()).PartialUpdate("123", user)
	assert.NotNil(t, err)
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userServiceMock := services.NewMockUserService(ctrl)
	userServiceMock.EXPECT().Delete("123").Return(nil)

	err := NewUserCase(getConfig(t), userServiceMock, configs.NewLog()).Delete("123")
	assert.Nil(t, err)
}

func TestDeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userServiceMock := services.NewMockUserService(ctrl)
	userServiceMock.EXPECT().Delete("123").Return(fmt.Errorf("123"))

	err := NewUserCase(getConfig(t), userServiceMock, configs.NewLog()).Delete("123")
	assert.NotNil(t, err)
}

func TestValidatePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pass, hasErr := helpers.HashPassword("32131")
	assert.Nil(t, hasErr)
	user := entity.User{ID: "123", Name: "Walisson Casonatto", Email: "wdcasonatto@gmail.com", Password: pass}

	userServiceMock := services.NewMockUserService(ctrl)
	userServiceMock.EXPECT().Find("123").Return(user, nil).Times(2)

	err := NewUserCase(getConfig(t), userServiceMock, configs.NewLog()).ValidatePassword("123", "32131")
	assert.Nil(t, err)

	err = NewUserCase(getConfig(t), userServiceMock, configs.NewLog()).ValidatePassword("123", "asdf")
	assert.NotNil(t, err)
}
