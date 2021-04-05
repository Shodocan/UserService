package requests

import (
	"testing"

	"github.com/Shodocan/UserService/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestUserValidation(t *testing.T) {
	req := SearchUserRequest{}

	assert.NotNil(t, req.Validate(), "must not be a valid request")
	req.Limit = 10
	assert.NotNil(t, req.Validate(), "must still not be a valid request")
	req.Page = 1
	assert.Nil(t, req.Validate(), "must be a valid request")
	req.Sort = []string{"name", "asdf"}
	assert.Nil(t, req.Validate(), "must be a valid request")
	req.Filters = []entity.UserFilter{{Field: "name", Operator: "=", Value: "123"}}
	assert.Nil(t, req.Validate(), "must be a valid request")
	req.Filters = []entity.UserFilter{{Field: "age", Operator: "~", Value: "12"}}
	assert.Nil(t, req.Validate(), "must be a valid request")
	req.Filters = []entity.UserFilter{{Field: "name", Operator: "*", Value: "123"}}
	assert.NotNil(t, req.Validate(), "must not be a valid request")
	req.Filters = []entity.UserFilter{{Field: "fest", Operator: "=", Value: "123"}}
	assert.NotNil(t, req.Validate(), "must not be a valid request")
}
