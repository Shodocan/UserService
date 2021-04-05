package mongo

import (
	"log"
	"os"
	"testing"

	"github.com/Shodocan/UserService/internal/configs"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// This test requires that MongoDb instance is Running and the testing data is loaded on the database

func checkMongo() bool {
	return os.Getenv("MONGODB_HOST") != ""
}

type TestEntity struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

func getConfig(t *testing.T) *configs.EnvVarConfig {
	config, err := configs.GetEnvConfig()
	if err != nil {
		t.Errorf("Error trying to get config %v", err)
		t.Fail()
	}
	return config
}

func TestMongoConnect(t *testing.T) {
	if !checkMongo() {
		t.Skip()
		return
	}
	mongo, err := NewDB(getConfig(t), configs.NewLog())
	assert.Nil(t, err, "Failed to start mongo", err)
	defer func() {
		err := mongo.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func TestPing(t *testing.T) {
	if !checkMongo() {
		t.Skip()
		return
	}
	mongo, err := NewDB(getConfig(t), configs.NewLog())
	assert.Nil(t, err, "Failed to start mongo", err)
	defer func() {
		err := mongo.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
	assert.Nil(t, mongo.Ping(), "Must ping")
}

func TestCreate(t *testing.T) {
	if !checkMongo() {
		t.Skip()
		return
	}
	mongo, err := NewDB(getConfig(t), configs.NewLog())
	assert.Nil(t, err, "Failed to start mongo", err)
	defer func() {
		err := mongo.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
	id, err := mongo.Create("teste", &TestEntity{Name: "Walisson"})
	assert.Nil(t, err, "Failed to save", err)
	assert.NotNil(t, id, "Failed to save", err)
	err = mongo.Delete("teste", id)
	assert.Nil(t, err, "Failed to delete", err)
}

func TestFind(t *testing.T) {
	if !checkMongo() {
		t.Skip()
		return
	}
	mongo, err := NewDB(getConfig(t), configs.NewLog())
	assert.Nil(t, err, "Failed to start mongo", err)
	defer func() {
		err := mongo.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
	id, err := mongo.Create("teste", &TestEntity{Name: "Walisson"})
	assert.Nil(t, err, "Failed to save", err)
	assert.NotNil(t, id, "Failed to save", err)

	testEntity := &TestEntity{}
	err = mongo.Find("teste", id, testEntity)
	assert.Nil(t, err, "Failed to find", err)
	assert.Equal(t, "Walisson", testEntity.Name, "name must be Walisson")
	err = mongo.Delete("teste", id)
	assert.Nil(t, err, "Failed to delete", err)
}

func TestUpdate(t *testing.T) {
	if !checkMongo() {
		t.Skip()
		return
	}
	mongo, err := NewDB(getConfig(t), configs.NewLog())
	assert.Nil(t, err, "Failed to start mongo", err)
	defer func() {
		err := mongo.Disconnect()
		if err != nil {
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	id, err := mongo.Create("teste", &TestEntity{Name: "Walisson"})
	assert.Nil(t, err, "Failed to save", err)
	assert.NotNil(t, id, "Failed to save", err)

	err = mongo.Update("teste", id, &TestEntity{Name: "Walisson Updated"})
	assert.Nil(t, err, "Failed to update", err)

	testEntity := &TestEntity{}
	err = mongo.Find("teste", id, testEntity)
	assert.Nil(t, err, "Failed to find", err)
	assert.Equal(t, "Walisson Updated", testEntity.Name, "name must be Walisson Updated")

	err = mongo.Delete("teste", id)
	assert.Nil(t, err, "Failed to delete", err)
}

func TestDelete(t *testing.T) {
	if !checkMongo() {
		t.Skip()
		return
	}
	mongo, err := NewDB(getConfig(t), configs.NewLog())
	assert.Nil(t, err, "Failed to start mongo", err)
	defer func() {
		err := mongo.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
	id, err := mongo.Create("teste", &TestEntity{Name: "Walisson"})
	assert.Nil(t, err, "Failed to save", err)
	assert.NotNil(t, id, "Failed to save", err)

	err = mongo.Delete("teste", id)
	assert.Nil(t, err, "Failed to delete", err)

	testEntity := &TestEntity{}
	err = mongo.Find("teste", id, testEntity)
	assert.NotNil(t, err, "Must fail to find", err)
}

func TestQuery(t *testing.T) {
	if !checkMongo() {
		t.Skip()
		return
	}
	mongo, err := NewDB(getConfig(t), configs.NewLog())
	assert.Nil(t, err, "Failed to start mongo", err)
	defer func() {
		err := mongo.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
	testCases := []TestEntity{
		{Name: "Walisson Casonatto"},
		{Name: "Bruce Lee"},
		{Name: "Luke Skywalker"},
		{Name: "Tyrion Lannister"},
		{Name: "Max Eisenhardt"},
		{Name: "Anakin Skywalker"},
	}

	for _, test := range testCases {
		id, err := mongo.Create("teste", &test)
		assert.Nil(t, err, "Failed to save", err)
		assert.NotNil(t, id, "Failed to save", err)
		defer func() {
			err = mongo.Delete("teste", id)
			assert.Nil(t, err, "Failed to delete", err)
		}()
	}

	res := []TestEntity{}
	err = mongo.Query("teste", bson.M{}, bson.M{}, 0, 6, &res)
	assert.Nil(t, err, "Should query")
	assert.Len(t, res, 6, "Must find 6 cases")

	res = []TestEntity{}
	err = mongo.Query("teste", bson.M{}, bson.M{}, 0, 5, &res)
	assert.Nil(t, err, "Should query")
	assert.Len(t, res, 5, "Must find 5 cases")

	res = []TestEntity{}
	err = mongo.Query("teste", bson.M{}, bson.M{"name": -1}, 0, 5, &res)
	assert.Nil(t, err, "Should query")
	assert.Len(t, res, 5, "Must find 5 cases")
	assert.Equal(t, res[0].Name, "Walisson Casonatto", "Must find 5 cases")

	res = []TestEntity{}
	err = mongo.Query("teste", bson.M{"name": "Walisson Casonatto"}, bson.M{}, 0, 5, &res)
	assert.Nil(t, err, "Should query")
	assert.Len(t, res, 1, "Must find 1 case")
	assert.Equal(t, res[0].Name, "Walisson Casonatto", "Must find 5 cases")

	res = []TestEntity{}
	err = mongo.Query("teste", bson.M{"name": bson.M{"$regex": "Skywalker"}}, primitive.D{{Key: "name", Value: 1}, {Key: "email", Value: 1}}, 0, 5, &res)
	assert.Nil(t, err, "Should query")
	assert.Len(t, res, 2, "Must find 2 cases")
	assert.Equal(t, res[0].Name, "Anakin Skywalker", "Must find 5 cases")
}

func TestTotal(t *testing.T) {
	if !checkMongo() {
		t.Skip()
		return
	}
	mongo, err := NewDB(getConfig(t), configs.NewLog())
	assert.Nil(t, err, "Failed to start mongo", err)
	defer func() {
		err := mongo.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
	testCases := []TestEntity{
		{Name: "Walisson Casonatto"},
		{Name: "Bruce Lee"},
		{Name: "Luke Skywalker"},
		{Name: "Tyrion Lannister"},
		{Name: "Max Eisenhardt"},
		{Name: "Anakin Skywalker"},
	}

	for _, test := range testCases {
		id, err := mongo.Create("teste", &test)
		assert.Nil(t, err, "Failed to save", err)
		assert.NotNil(t, id, "Failed to save", err)
		defer func() {
			err = mongo.Delete("teste", id)
			assert.Nil(t, err, "Failed to delete", err)
		}()
	}

	var total int
	total, err = mongo.Total("teste", bson.M{})
	assert.Nil(t, err, "Should query")
	assert.Equal(t, total, 6, "Must find 6 cases")

	total, err = mongo.Total("teste", bson.M{"name": "Walisson Casonatto"})
	assert.Nil(t, err, "Should query")
	assert.Equal(t, total, 1, "Must find 1 cases")

	total, err = mongo.Total("teste", bson.M{"name": bson.M{"$regex": "Skywalker"}})
	assert.Nil(t, err, "Should query")
	assert.Equal(t, total, 2, "Must find 1 cases")
}
