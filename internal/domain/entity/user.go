package entity

import (
	"net/http"

	"github.com/Shodocan/UserService/internal/configs/engine"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

func (u User) Validate(creation bool) error {
	validationErrors := map[string]interface{}{}
	if u.Email == "" {
		validationErrors["email"] = "Email is required"
	}
	if u.Name == "" {
		validationErrors["name"] = "Name is required"
	}
	if u.Password == "" && creation {
		validationErrors["password"] = "Password is required"
	}
	if len(validationErrors) > 0 {
		return engine.NewGenericError(http.StatusBadRequest, "Invalid User").ExtraData(validationErrors)
	}
	return nil
}

type UserFied string

const (
	ID      UserFied = "_id"
	Name    UserFied = "name"
	Age     UserFied = "age"
	Email   UserFied = "email"
	Address UserFied = "address"
)

type UserFilter struct {
	Field    UserFied `enums:"name,age,email,address"`
	Value    interface{}
	Operator FilterOperation `enums:"=,~"`
}

func (f UserFilter) Valid() bool {
	switch f.Field {
	case Name, Age, Address, Email:
		break
	default:
		return false
	}

	switch f.Operator {
	case Equal, Like:
		break
	default:
		return false
	}

	return true
}
