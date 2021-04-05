package mongoredis

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shodocan/UserService/internal/configs"
	"github.com/Shodocan/UserService/internal/database"
	"github.com/Shodocan/UserService/internal/database/mongo"
	"github.com/Shodocan/UserService/internal/database/redis"
)

type DB struct {
	mongo  database.MongoDB
	redis  database.RedisDB
	logger *log.Logger
}

func NewDB(config *configs.EnvVarConfig, logger *log.Logger) (database.MongoDB, error) {
	mongodb, err := mongo.NewDB(config, logger)
	if err != nil {
		return nil, err
	}

	redisdb, err := redis.NewDB(config, logger)
	if err != nil {
		logger.Printf("Failed to connect redis :%v", err)
	}

	return &DB{mongo: mongodb, redis: redisdb, logger: logger}, nil
}

func (db *DB) Disconnect() error {
	err := db.mongo.Disconnect()
	redisErr := db.redis.Disconnect()

	if err != nil || redisErr != nil {
		return fmt.Errorf("error disconecting redis: %v, mongo: %v", redisErr, err)
	}
	return nil
}

func (db *DB) Ping() error {
	err := db.redis.Ping()
	if err != nil {
		db.logger.Printf("failed to ping redis %v", err)
	}
	return db.mongo.Ping()
}

func (db DB) Create(namespace string, data interface{}) (string, error) {
	id, err := db.mongo.Create(namespace, data)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (db DB) Find(namespace, idStr string, dst interface{}) error {
	err := db.redis.Get(idStr, dst)
	if err != nil {
		err := db.mongo.Find(namespace, idStr, dst)
		if err != nil {
			return err
		}
		_, err = db.redis.Set(idStr, dst)
		if err != nil {
			db.logger.Println("failed to cache on Redis")
		}
	}
	return nil
}

func (db DB) Query(namespace string, filters, sort interface{}, offset, limit int, dst interface{}) error {
	hash := db.QueryHash(namespace, filters, sort, offset, limit)

	err := db.redis.Get(hash, dst)
	if err != nil {
		err = db.mongo.Query(namespace, filters, sort, offset, limit, dst)
		if err != nil {
			return err
		}
		_, err = db.redis.Set(hash, dst)
		if err != nil {
			db.logger.Println("failed to cache on Redis")
		}
	}

	return nil
}

func (db DB) Total(namespace string, filters interface{}) (int, error) {
	hash := db.QueryHash(namespace, filters, nil, -1, -1)

	var total int
	err := db.redis.Get(hash, &total)
	if err != nil {
		total, err = db.mongo.Total(namespace, filters)
		if err != nil {
			return total, err
		}
		_, err = db.redis.Set(hash, total)
		if err != nil {
			db.logger.Println("failed to cache on Redis")
		}
	}

	return total, nil
}

func (db DB) Update(namespace, idStr string, data interface{}) error {
	err := db.mongo.Update(namespace, idStr, data)
	if err != nil {
		return err
	}
	err = db.redis.Delete(idStr)
	if err != nil {
		db.logger.Println("failed to cache on Redis")
	}
	return nil
}

func (db DB) Delete(namespace, idStr string) error {
	err := db.mongo.Delete(namespace, idStr)
	redisErr := db.redis.Delete(idStr)

	if db.redis.Ping() != nil {
		// if redis is unvaliable must ignore delete errors
		redisErr = nil
	}

	if err != nil || redisErr != nil {
		return fmt.Errorf("error deleting redis: %v, mongo:%v", redisErr, err)
	}
	return nil
}

func (db DB) QueryHash(namespace string, filters, sort interface{}, offset, limit int) string {
	q := Query{
		Namespace: namespace,
		Filters:   filters,
		Sort:      sort,
		Offset:    offset,
		Limit:     limit,
	}
	jsonData, err := json.Marshal(q)
	if err != nil {
		db.logger.Println(err)
	}
	return base64.RawStdEncoding.EncodeToString(jsonData)
}

type Query struct {
	Namespace string
	Filters   interface{}
	Sort      interface{}
	Offset    int
	Limit     int
}
