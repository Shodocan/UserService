package mongo

import (
	"context"
	"time"

	"github.com/Shodocan/UserService/internal/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func Migrate(config *configs.EnvVarConfig) error {
	dbconfig := config.GetMongoConnectionString()
	client, err := mongo.NewClient(options.Client().ApplyURI(dbconfig), &options.ClientOptions{Auth: &options.Credential{
		Username: config.MongoDBAdminUsername,
		Password: config.MongodbAdminPassword,
	}})
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	_, err = client.Database(config.MongoDBDatabase).Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"email": 1, // index in ascending order
		}, Options: options.Index().SetUnique(true),
	})
	return err
}
