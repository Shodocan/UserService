package services

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/Shodocan/UserService/internal/configs"
	"github.com/Shodocan/UserService/internal/database"
	"github.com/Shodocan/UserService/internal/database/mongo"
	"github.com/Shodocan/UserService/internal/database/mongoredis"
	"github.com/Shodocan/UserService/internal/domain/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func checkMongo() bool {
	return os.Getenv("MONGODB_HOST") != ""
}

func checkRedis() bool {
	return os.Getenv("REDISDB_HOST") != ""
}

func getConfig(t *testing.T) *configs.EnvVarConfig {
	config, err := configs.GetEnvConfig()
	if err != nil {
		t.Log(err)
		return config
	}
	return config
}

func getMongoRedis(
	config *configs.EnvVarConfig,
	ctrl *gomock.Controller,
	configMock func(*database.MockMongoDB) *database.MockMongoDB,
) (database.MongoDB, error) {
	if checkMongo() && checkRedis() {
		return mongoredis.NewDB(config, configs.NewLog())
	} else if checkMongo() {
		return mongo.NewDB(config, configs.NewLog())
	} else {
		return configMock(database.NewMockMongoDB(ctrl)), nil
	}
}

func TestMongoRepositoryCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cfg := getConfig(t)

	user := entity.User{
		Name:     "Walisson Casonatto",
		Age:      28,
		Email:    "wdcasonatto@gmail.com",
		Password: "3020102030",
		Address:  "Rua das ...",
	}

	db, err := getMongoRedis(cfg, ctrl, func(mmd *database.MockMongoDB) *database.MockMongoDB {
		userID := primitive.NewObjectID().Hex()
		gomock.InOrder(
			mmd.EXPECT().Create("users", gomock.Any()).Return(userID, nil),
			mmd.EXPECT().Find("users", userID, gomock.Any()).DoAndReturn(func(namespace, id string, dst interface{}) error {
				if userID == id {
					dstVal := reflect.ValueOf(dst).Elem()
					id, _ := primitive.ObjectIDFromHex(id)
					dstVal.FieldByName("ID").Set(reflect.ValueOf(id))
					dstVal.FieldByName("Name").SetString(user.Name)
					dstVal.FieldByName("Age").SetInt(int64(user.Age))
					dstVal.FieldByName("Email").SetString(user.Email)
					dstVal.FieldByName("Password").SetString(user.Password)
					dstVal.FieldByName("Address").SetString(user.Address)
				}
				return nil
			}),
			mmd.EXPECT().Delete("users", userID).Return(nil),
			mmd.EXPECT().Disconnect().Return(nil),
		)
		return mmd
	})
	defer func() {
		err := db.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
	assert.Nil(t, err, "Should return a mongo instance")

	repo := NewUserServiceMongo(db, cfg)

	createdUser, err := repo.Create(user)
	assert.Nil(t, err, "Should return a new user")
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Age, createdUser.Age)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Password, createdUser.Password)
	assert.Equal(t, user.Address, createdUser.Address)

	err = repo.Delete(createdUser.ID)
	assert.Nil(t, err, "Should delete user")
}

func TestMongoRepositoryUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cfg := getConfig(t)

	user := entity.User{
		Name:     "Walisson Casonatto",
		Age:      28,
		Email:    "wdcasonatto@gmail.com",
		Password: "3020102030",
		Address:  "Rua das ...",
	}

	updateUser := entity.User{
		Name:     "Walisson De Deus Casonatto",
		Age:      30,
		Email:    "wdcasonatto@gmail.com",
		Password: "",
		Address:  "",
	}

	db, err := getMongoRedis(cfg, ctrl, func(mmd *database.MockMongoDB) *database.MockMongoDB {
		userID := primitive.NewObjectID().Hex()
		gomock.InOrder(
			mmd.EXPECT().Create("users", gomock.Any()).Return(userID, nil),
			mmd.EXPECT().Find("users", userID, gomock.Any()).DoAndReturn(func(namespace, id string, dst interface{}) error {
				if userID == id {
					dstVal := reflect.ValueOf(dst).Elem()
					id, _ := primitive.ObjectIDFromHex(id)
					dstVal.FieldByName("ID").Set(reflect.ValueOf(id))
					dstVal.FieldByName("Name").SetString(user.Name)
					dstVal.FieldByName("Age").SetInt(int64(user.Age))
					dstVal.FieldByName("Email").SetString(user.Email)
					dstVal.FieldByName("Password").SetString(user.Password)
					dstVal.FieldByName("Address").SetString(user.Address)
				}
				return nil
			}),
			mmd.EXPECT().Update("users", userID, gomock.Any()).Return(nil),
			mmd.EXPECT().Find("users", userID, gomock.Any()).DoAndReturn(func(namespace, id string, dst interface{}) error {
				if userID == id {
					dstVal := reflect.ValueOf(dst).Elem()
					id, _ := primitive.ObjectIDFromHex(id)
					dstVal.FieldByName("ID").Set(reflect.ValueOf(id))
					dstVal.FieldByName("Name").SetString(updateUser.Name)
					dstVal.FieldByName("Age").SetInt(int64(updateUser.Age))
					dstVal.FieldByName("Email").SetString(updateUser.Email)
					dstVal.FieldByName("Password").SetString(user.Password)
					dstVal.FieldByName("Address").SetString(updateUser.Address)
				}
				return nil
			}),
			mmd.EXPECT().Update("users", userID, gomock.Any()).Return(nil),
			mmd.EXPECT().Find("users", userID, gomock.Any()).DoAndReturn(func(namespace, id string, dst interface{}) error {
				if userID == id {
					dstVal := reflect.ValueOf(dst).Elem()
					id, _ := primitive.ObjectIDFromHex(id)
					dstVal.FieldByName("ID").Set(reflect.ValueOf(id))
					dstVal.FieldByName("Name").SetString(updateUser.Name)
					dstVal.FieldByName("Age").SetInt(int64(updateUser.Age))
					dstVal.FieldByName("Email").SetString(updateUser.Email)
					dstVal.FieldByName("Password").SetString(updateUser.Password)
					dstVal.FieldByName("Address").SetString(updateUser.Address)
				}
				return nil
			}),
			mmd.EXPECT().Delete("users", userID).Return(nil),
			mmd.EXPECT().Disconnect().Return(nil),
		)
		return mmd
	})
	defer func() {
		err := db.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
	assert.Nil(t, err, "Should return a mongo instance")

	repo := NewUserServiceMongo(db, cfg)

	createdUser, err := repo.Create(user)
	assert.Nil(t, err, "Should return a new user")
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Age, createdUser.Age)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Password, createdUser.Password)
	assert.Equal(t, user.Address, createdUser.Address)

	updateUser.ID = createdUser.ID
	updatedUser, err := repo.Update(createdUser.ID, updateUser)
	assert.Nil(t, err, "Should return a updated user")
	assert.Equal(t, createdUser.ID, updatedUser.ID)
	assert.Equal(t, updateUser.Name, updatedUser.Name)
	assert.Equal(t, updateUser.Age, updatedUser.Age)
	assert.Equal(t, updateUser.Email, updatedUser.Email)
	assert.Equal(t, user.Password, updatedUser.Password, "must not change password")
	assert.Equal(t, updateUser.Address, updatedUser.Address)

	updateUser.Password = "333333333"
	updatedUser, err = repo.Update(createdUser.ID, updateUser)
	assert.Nil(t, err, "Should return a updated user")
	assert.Equal(t, updateUser.Name, updatedUser.Name)
	assert.Equal(t, updateUser.Age, updatedUser.Age)
	assert.Equal(t, updateUser.Email, updatedUser.Email)
	assert.Equal(t, updateUser.Password, updatedUser.Password, "must change password")
	assert.Equal(t, updateUser.Address, updatedUser.Address)

	err = repo.Delete(updatedUser.ID)
	assert.Nil(t, err, "Should delete user")
}

func TestMongoRepositoryPartialUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cfg := getConfig(t)

	user := entity.User{
		Name:     "Walisson Casonatto",
		Age:      28,
		Email:    "wdcasonatto@gmail.com",
		Password: "3020102030",
		Address:  "Rua das ...",
	}

	updateUser := entity.User{
		Name:  "Walisson De Deus Casonatto",
		Email: "wdcasonatto@gmail.com",
	}

	db, err := getMongoRedis(cfg, ctrl, func(mmd *database.MockMongoDB) *database.MockMongoDB {
		userID := primitive.NewObjectID().Hex()
		gomock.InOrder(
			mmd.EXPECT().Create("users", gomock.Any()).Return(userID, nil),
			mmd.EXPECT().Find("users", userID, gomock.Any()).DoAndReturn(func(namespace, id string, dst interface{}) error {
				if userID == id {
					dstVal := reflect.ValueOf(dst).Elem()
					id, _ := primitive.ObjectIDFromHex(id)
					dstVal.FieldByName("ID").Set(reflect.ValueOf(id))
					dstVal.FieldByName("Name").SetString(user.Name)
					dstVal.FieldByName("Age").SetInt(int64(user.Age))
					dstVal.FieldByName("Email").SetString(user.Email)
					dstVal.FieldByName("Password").SetString(user.Password)
					dstVal.FieldByName("Address").SetString(user.Address)
				}
				return nil
			}),
			mmd.EXPECT().Update("users", userID, gomock.Any()).Return(nil),
			mmd.EXPECT().Find("users", userID, gomock.Any()).DoAndReturn(func(namespace, id string, dst interface{}) error {
				if userID == id {
					dstVal := reflect.ValueOf(dst).Elem()
					id, _ := primitive.ObjectIDFromHex(id)
					dstVal.FieldByName("ID").Set(reflect.ValueOf(id))
					dstVal.FieldByName("Name").SetString(updateUser.Name)
					dstVal.FieldByName("Age").SetInt(int64(user.Age))
					dstVal.FieldByName("Email").SetString(updateUser.Email)
					dstVal.FieldByName("Password").SetString(user.Password)
					dstVal.FieldByName("Address").SetString(user.Address)
				}
				return nil
			}),
			mmd.EXPECT().Delete("users", userID).Return(nil),
			mmd.EXPECT().Disconnect().Return(nil),
		)
		return mmd
	})
	defer func() {
		err := db.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
	assert.Nil(t, err, "Should return a mongo instance")

	repo := NewUserServiceMongo(db, cfg)

	createdUser, err := repo.Create(user)
	assert.Nil(t, err, "Should return a new user")
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Age, createdUser.Age)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Password, createdUser.Password)
	assert.Equal(t, user.Address, createdUser.Address)

	updateUser.ID = createdUser.ID
	updatedUser, err := repo.PartialUpdate(createdUser.ID, updateUser)
	assert.Nil(t, err, "Should return a updated user")
	assert.Equal(t, createdUser.ID, updatedUser.ID)
	assert.Equal(t, updateUser.Name, updatedUser.Name)
	assert.Equal(t, user.Age, updatedUser.Age, "must not change password")
	assert.Equal(t, updateUser.Email, updatedUser.Email)
	assert.Equal(t, user.Password, updatedUser.Password, "must not change password")
	assert.Equal(t, user.Address, updatedUser.Address, "must not change Address")

	err = repo.Delete(updatedUser.ID)
	assert.Nil(t, err, "Should delete user")
}

func TestMongoRepositoryQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cfg := getConfig(t)

	users := []*entity.User{
		{
			Name:     "Walisson Casonatto",
			Age:      28,
			Email:    "wdcasonatto@gmail.com",
			Password: "3020102030",
			Address:  "Rua das ...",
		},
		{
			Name:     "Samara Casonatto",
			Age:      24,
			Email:    "samaracasonatto@gmail.com",
			Password: "5050505050",
			Address:  "Rua das ...",
		},
		{
			Name: "Luke Skywalker",
			Age:  24,
		},
		{
			Name: "Anakin Skywalker",
			Age:  48,
		},
	}

	db, err := getMongoRedis(cfg, ctrl, func(mmd *database.MockMongoDB) *database.MockMongoDB {
		t.SkipNow()
		return mmd
	})
	defer func() {
		err := db.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
	assert.Nil(t, err, "Should return a mongo instance")

	repo := NewUserServiceMongo(db, cfg)
	createdUsers := []entity.User{}

	for _, user := range users {
		createdUser, err := repo.Create(*user)
		assert.Nil(t, err, "Should return a new user")
		assert.Equal(t, user.Name, createdUser.Name)
		assert.Equal(t, user.Age, createdUser.Age)
		assert.Equal(t, user.Email, createdUser.Email)
		assert.Equal(t, user.Password, createdUser.Password)
		assert.Equal(t, user.Address, createdUser.Address)
		createdUsers = append(createdUsers, createdUser)
	}

	foundUsers, pagination, err := repo.Query([]entity.UserFilter{{Field: "name", Operator: entity.Like, Value: "Skywalker"}}, []string{"-name", "age"}, 1, 100)
	assert.Nil(t, err, "Must query")
	assert.Len(t, foundUsers, 2, "Should find all skywalkers")
	assert.Equal(t, "Luke Skywalker", foundUsers[0].Name, "Luke should be the first")
	assert.NotNil(t, pagination)

	foundUsers, pagination, err = repo.Query([]entity.UserFilter{}, []string{"name", "age"}, 1, 2)
	assert.Nil(t, err, "Must query")
	assert.Len(t, foundUsers, 2, "Should find 2")
	assert.Equal(t, "Anakin Skywalker", foundUsers[0].Name, "Anakin should be the first")
	assert.Equal(t, "Luke Skywalker", foundUsers[1].Name, "Luke should be the second")
	assert.NotNil(t, pagination)

	foundUsers, pagination, err = repo.Query([]entity.UserFilter{}, []string{"name", "age"}, 2, 2)
	assert.Nil(t, err, "Must query")
	assert.Len(t, foundUsers, 2, "Should find 2")
	assert.Equal(t, "Samara Casonatto", foundUsers[0].Name, "Samara should be the first")
	assert.Equal(t, "Walisson Casonatto", foundUsers[1].Name, "Walisson should be the second")
	assert.NotNil(t, pagination)

	foundUsers, pagination, err = repo.Query([]entity.UserFilter{}, []string{"name", "age"}, 3, 2)
	assert.Nil(t, err, "Must query")
	assert.Len(t, foundUsers, 0, "Should find 0")
	assert.NotNil(t, pagination)

	for _, user := range createdUsers {
		err = repo.Delete(user.ID)
		assert.Nil(t, err, "Should delete user")
	}
}
