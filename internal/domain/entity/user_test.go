package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserFilterValidation(t *testing.T) {
	filter1 := UserFilter{Field: "name", Operator: "=", Value: "123"}
	filter2 := UserFilter{Field: "age", Operator: "=", Value: "123"}
	filter3 := UserFilter{Field: "address", Operator: "=", Value: "123"}
	filter4 := UserFilter{Field: "email", Operator: "=", Value: "123"}
	filter5 := UserFilter{Field: "password", Operator: "=", Value: "123"}
	filter6 := UserFilter{Field: "name", Operator: "~", Value: "123"}
	filter7 := UserFilter{Field: "name", Operator: "*", Value: "123"}
	filter8 := UserFilter{Field: "name", Operator: "-", Value: "123"}
	filter9 := UserFilter{Field: "asdf", Operator: "=", Value: "123"}

	assert.True(t, filter1.Valid())
	assert.True(t, filter2.Valid())
	assert.True(t, filter3.Valid())
	assert.True(t, filter4.Valid())
	assert.False(t, filter5.Valid())
	assert.True(t, filter6.Valid())
	assert.False(t, filter7.Valid())
	assert.False(t, filter8.Valid())
	assert.False(t, filter9.Valid())
}

func TestUserValidation(t *testing.T) {
	user1 := User{Name: "Walisson Casonatto", Email: "wdcasonatto@gmail.com", Password: "123"}
	user2 := User{Name: "Walisson Casonatto", Email: "wdcasonatto@gmail.com"}
	user3 := User{Name: "Walisson Casonatto", Password: "123"}
	user4 := User{Email: "wdcasonatto@gmail.com", Password: "123"}

	assert.Nil(t, user1.Validate(true), "Should be valid to create")
	assert.Nil(t, user1.Validate(false), "Should be valid to update")
	assert.NotNil(t, user2.Validate(true), "Should not be valid to create")
	assert.Nil(t, user2.Validate(false), "Should be valid to update")
	assert.NotNil(t, user3.Validate(true), "Should not be valid to create")
	assert.NotNil(t, user3.Validate(false), "Should not be valid to update")
	assert.NotNil(t, user4.Validate(true), "Should not be valid to create")
	assert.NotNil(t, user4.Validate(false), "Should not be valid to update")
}
