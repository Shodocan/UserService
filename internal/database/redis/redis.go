package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/KromDaniel/rejonson"
	"github.com/Shodocan/UserService/internal/configs"
	"github.com/Shodocan/UserService/internal/database"
	"github.com/go-redis/redis"
)

type DB struct {
	cache   time.Duration
	timeout time.Duration
	config  *configs.EnvVarConfig
	logger  *log.Logger
	_client *rejonson.Client
}

func NewDB(config *configs.EnvVarConfig, logger *log.Logger) (database.RedisDB, error) {
	cache, err := time.ParseDuration(config.RedisDBCacheDuration)
	if err != nil {
		logger.Printf("Invalid Redis cache duration %s, using 30s", config.RedisDBCacheDuration)
		cache = 30 * time.Second
	}
	timeout, err := time.ParseDuration(config.RedisDBDefaultTimeout)
	if err != nil {
		logger.Printf("Invalid Redis cache duration %s, using 3s", config.RedisDBDefaultTimeout)
		cache = 3 * time.Second
	}
	db := &DB{
		cache:   cache,
		logger:  logger,
		config:  config,
		timeout: timeout,
	}

	err = db.init()
	return db, err
}

func (db *DB) init() error {
	var err error
	dbs, err := strconv.Atoi(db.config.RedisDBDatabase)
	if err != nil {
		db.logger.Printf("Invalid Redis Database %s, using 1", db.config.RedisDBDatabase)
		dbs = 1
	}
	goRedisClient := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%s", db.config.RedisDBHost, db.config.RedisDBPort),
		Password:    db.config.RedisDBPassword,
		MaxRetries:  5,
		DialTimeout: db.timeout,
		DB:          dbs,
	})
	client := rejonson.ExtendClient(goRedisClient)
	db._client = client
	return err
}

func (db *DB) Disconnect() error {
	return db._client.Close()
}

func (db *DB) Ping() error {
	return db._client.Ping().Err()
}

func (db DB) Set(idStr string, data interface{}) (string, error) {
	js, _ := json.Marshal(data)
	_, err := db._client.JsonSet(idStr, ".", string(js)).Result()
	db._client.Expire(idStr, db.cache)
	return idStr, err
}

func (db DB) Get(idStr string, dst interface{}) error {
	jsonString, err := db._client.JsonGet(idStr).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonString), dst)
}

func (db DB) Delete(idStr string) error {
	sts := db._client.Del(idStr)
	return sts.Err()
}
