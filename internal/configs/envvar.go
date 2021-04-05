package configs

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type EnvVarConfig struct {
	Port                  string `envconfig:"port" default:"8080"`
	APIToken              string `envconfig:"api_token" default:"e81384e6-2b68-4d40-b19e-dd585132baa9"`
	AllowOrigins          string `envconfig:"allowed_origins" default:"localhost"`
	MongoDBHost           string `envconfig:"mongodb_host" default:"mongodb"`
	MongoDBPort           string `envconfig:"mongodb_port" default:"27017"`
	MongoDBDatabase       string `envconfig:"mongodb_database" default:"user"`
	MongoDBDefaultTimeout string `envconfig:"mongodb_default_timeout" default:"10s"`
	MongoDBAdminUsername  string `envconfig:"mongodb_adminusername" required:"true"`
	MongodbAdminPassword  string `envconfig:"mongodb_adminpassword" required:"true"`
	RedisDBActive         string `envconfig:"redisdb_active" default:"false"`
	RedisDBHost           string `envconfig:"redisdb_host" default:""`
	RedisDBPort           string `envconfig:"redisdb_port" default:"6379"`
	RedisDBPassword       string `envconfig:"redisdb_password" default:"root"`
	RedisDBDatabase       string `envconfig:"redisdb_database" default:"1"`
	RedisDBCacheDuration  string `envconfig:"redisdb_cache_duration" default:"30s"`
	RedisDBDefaultTimeout string `envconfig:"redisdb_default_timeout" default:"3s"`
}

func GetEnvConfig() (*EnvVarConfig, error) {
	var cfg EnvVarConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (config *EnvVarConfig) GetMongoConnectionString() string {
	return fmt.Sprintf("mongodb://%s:%s", config.MongoDBHost, config.MongoDBPort)
}
