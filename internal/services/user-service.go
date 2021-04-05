package services

import (
	"strings"

	"github.com/Shodocan/UserService/internal/configs"
	"github.com/Shodocan/UserService/internal/configs/engine"
	"github.com/Shodocan/UserService/internal/database"
	"github.com/Shodocan/UserService/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func NewUserServiceMongo(db database.MongoDB, config *configs.EnvVarConfig) UserService {
	return &UserServiceMongo{_db: db, config: config}
}

type UserServiceMongo struct {
	_db    database.MongoDB
	config *configs.EnvVarConfig
}

func MapDBUser(user entity.User) *DBUser {
	dbUser := &DBUser{}
	if user.ID != "" {
		id, _ := primitive.ObjectIDFromHex(user.ID)
		dbUser.ID = id
	}
	dbUser.Address = user.Address
	dbUser.Password = user.Password
	dbUser.Age = user.Age
	dbUser.Name = user.Name
	dbUser.Email = user.Email
	return dbUser
}

func MapPartialUser(user entity.User) *DBUserPartial {
	dbUser := &DBUserPartial{}
	if user.ID != "" {
		id, _ := primitive.ObjectIDFromHex(user.ID)
		dbUser.ID = id
	}
	dbUser.Address = user.Address
	dbUser.Password = user.Password
	dbUser.Age = user.Age
	dbUser.Name = user.Name
	dbUser.Email = user.Email
	return dbUser
}

type DBUserPartial struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name,omitempty"`
	Age      int                `json:"age" bson:"age,omitempty"`
	Email    string             `json:"email" bson:"email,omitempty"`
	Password string             `json:"password" bson:"password,omitempty"`
	Address  string             `json:"address" bson:"address,omitempty"`
}

type DBUserList []DBUser

type DBUser struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Age      int                `json:"age" bson:"age"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password,omitempty"`
	Address  string             `json:"address" bson:"address"`
}

func (dbUser *DBUser) ToUser() entity.User {
	return entity.User{
		ID:       dbUser.ID.Hex(),
		Name:     dbUser.Name,
		Age:      dbUser.Age,
		Email:    dbUser.Email,
		Password: dbUser.Password,
		Address:  dbUser.Address,
	}
}

func (list DBUserList) ToUserList() []entity.User {
	result := []entity.User{}
	for _, usr := range list {
		result = append(result, usr.ToUser())
	}
	return result
}

func (repo UserServiceMongo) operationResolver(op entity.FilterOperation) string {
	switch op {
	case entity.Equal:
		return "$eq"
	default:
		return "$regex"
	}
}

func (repo UserServiceMongo) Query(filters []entity.UserFilter, sortList []string, page, limit int) ([]entity.User, engine.Pagination, error) {
	queryFilter := bson.M{}
	for _, filter := range filters {
		if filter.Valid() {
			queryFilter[string(filter.Field)] = bson.M{repo.operationResolver(filter.Operator): filter.Value}
		}
	}

	sort := primitive.D{}
	for _, rule := range sortList {
		if strings.Contains(rule, "-") {
			sort = append(sort, primitive.E{Key: strings.ReplaceAll(rule, "-", ""), Value: -1})
		} else {
			sort = append(sort, primitive.E{Key: rule, Value: 1})
		}
	}

	total, err := repo._db.Total("users", queryFilter)
	if err != nil {
		return nil, engine.Pagination{}, err
	}

	dbUsers := DBUserList{}
	err = repo._db.Query("users", queryFilter, sort, (page*limit)-limit, limit, &dbUsers)
	if err != nil {
		return nil, engine.Pagination{}, err
	}
	return dbUsers.ToUserList(), engine.NewPagination(total, len(dbUsers), limit), err
}

func (repo UserServiceMongo) Create(data entity.User) (entity.User, error) {
	dbUser := MapDBUser(data)
	id, err := repo._db.Create("users", dbUser)
	if err != nil {
		return entity.User{}, err
	}
	return repo.Find(id)
}

func (repo UserServiceMongo) Find(id string) (entity.User, error) {
	mongoUser := &DBUser{}
	err := repo._db.Find("users", id, mongoUser)
	return mongoUser.ToUser(), err
}

func (repo UserServiceMongo) Update(id string, user entity.User) (entity.User, error) {
	mongoUser := MapDBUser(user)
	err := repo._db.Update("users", id, mongoUser)
	if err != nil {
		return entity.User{}, err
	}
	return repo.Find(id)
}

func (repo UserServiceMongo) PartialUpdate(id string, user entity.User) (entity.User, error) {
	mongoUser := MapPartialUser(user)
	err := repo._db.Update("users", id, mongoUser)
	if err != nil {
		return entity.User{}, err
	}
	return repo.Find(id)
}

func (repo UserServiceMongo) Delete(id string) error {
	return repo._db.Delete("users", id)
}
