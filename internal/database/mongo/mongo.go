package mongo

import (
	"context"
	"log"
	"time"

	"github.com/Shodocan/UserService/internal/configs"
	"github.com/Shodocan/UserService/internal/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

type DB struct {
	timeout time.Duration
	config  *configs.EnvVarConfig
	logger  *log.Logger
	_client *mongo.Client
}

func NewDB(config *configs.EnvVarConfig, logger *log.Logger) (database.MongoDB, error) {
	timeout, err := time.ParseDuration(config.MongoDBDefaultTimeout)
	if err != nil {
		logger.Printf("Invalid MongoDB timeout %s, using 10s", config.MongoDBDefaultTimeout)
		timeout = 10 * time.Second
	}
	db := &DB{
		timeout: timeout,
		logger:  logger,
		config:  config,
	}

	err = db.init()
	return db, err
}

func (db *DB) init() error {
	var err error
	config := db.config.GetMongoConnectionString()
	db._client, err = mongo.NewClient(options.Client().ApplyURI(config), &options.ClientOptions{Auth: &options.Credential{
		Username: db.config.MongoDBAdminUsername,
		Password: db.config.MongodbAdminPassword,
	}})
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()
	db.logger.Printf("Conectado ao mongo %s %s", db.config.MongoDBHost, db.config.MongoDBDatabase)
	return db._client.Connect(ctx)
}

func (db *DB) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()
	return db._client.Disconnect(ctx)
}

func (db *DB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()
	return db._client.Ping(ctx, readpref.Primary())
}

func (db DB) Create(namespace string, data interface{}) (string, error) {
	collection := db._client.Database(db.config.MongoDBDatabase).Collection(namespace)

	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return "", err
	}
	id := res.InsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}

func (db DB) Find(namespace, idStr string, dst interface{}) error {
	collection := db._client.Database(db.config.MongoDBDatabase).Collection(namespace)
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return err
	}

	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(dst)
	return err
}

func (db DB) Query(namespace string, filters, sort interface{}, offset, limit int, dst interface{}) error {
	collection := db._client.Database(db.config.MongoDBDatabase).Collection(namespace)

	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(sort)

	cur, err := collection.Find(ctx, filters, findOptions)
	if err != nil {
		return err
	}
	return cur.All(ctx, dst)
}

func (db DB) Total(namespace string, filters interface{}) (int, error) {
	collection := db._client.Database(db.config.MongoDBDatabase).Collection(namespace)

	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	total, err := collection.CountDocuments(ctx, filters)
	if err != nil {
		return 0, err
	}
	return int(total), nil
}

func (db DB) Update(namespace, idStr string, data interface{}) error {
	collection := db._client.Database(db.config.MongoDBDatabase).Collection(namespace)

	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return err
	}

	_, err = collection.UpdateByID(ctx, id, bson.M{"$set": data})
	return err
}

func (db DB) Delete(namespace, idStr string) error {
	collection := db._client.Database(db.config.MongoDBDatabase).Collection(namespace)
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
