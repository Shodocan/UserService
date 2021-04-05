package redis

import (
	"log"
	"os"
	"testing"

	"github.com/Shodocan/UserService/internal/configs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// This test requires that redisDb instance is Running and the testing data is loaded on the database

func checkRedis() bool {
	return os.Getenv("REDISDB_HOST") != ""
}

type TestEntity struct {
	ID   string `json:"_id,omitempty"`
	Name string `json:"name"`
}

func getConfig(t *testing.T) *configs.EnvVarConfig {
	config, err := configs.GetEnvConfig()
	if err != nil {
		t.Errorf("Error trying to get config %v", err)
		t.Fail()
	}
	return config
}

func TestRedisConnect(t *testing.T) {
	if !checkRedis() {
		t.Skip()
		return
	}
	redis, err := NewDB(getConfig(t), configs.NewLog())
	assert.Nil(t, err, "Failed to start redis", err)
	defer func() {
		err := redis.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func TestPing(t *testing.T) {
	if !checkRedis() {
		t.Skip()
		return
	}
	redis, err := NewDB(getConfig(t), configs.NewLog())
	assert.Nil(t, err, "Failed to start redis", err)
	defer func() {
		err := redis.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
	assert.Nil(t, redis.Ping(), "Must ping")
}

func TestSet(t *testing.T) {
	if !checkRedis() {
		t.Skip()
		return
	}
	redis, err := NewDB(getConfig(t), configs.NewLog())
	assert.Nil(t, err, "Failed to start redis", err)
	defer func() {
		err := redis.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
	id, err := redis.Set(uuid.New().String(), &TestEntity{Name: "Walisson"})
	assert.Nil(t, err, "Failed to save", err)
	assert.NotNil(t, id, "Failed to save", err)
	err = redis.Delete(id)
	assert.Nil(t, err, "Failed to delete", err)
}

func TestFind(t *testing.T) {
	if !checkRedis() {
		t.Skip()
		return
	}
	redis, err := NewDB(getConfig(t), configs.NewLog())
	assert.Nil(t, err, "Failed to start redis", err)
	defer func() {
		err := redis.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
	id, err := redis.Set(uuid.NewString(), &TestEntity{Name: "Walisson"})
	assert.Nil(t, err, "Failed to save", err)
	assert.NotNil(t, id, "Failed to save", err)

	testEntity := &TestEntity{}
	err = redis.Get(id, testEntity)
	assert.Nil(t, err, "Failed to find", err)
	assert.Equal(t, "Walisson", testEntity.Name, "name must be Walisson")
	err = redis.Delete(id)
	assert.Nil(t, err, "Failed to delete", err)
}
