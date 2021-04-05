package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserPasswort(t *testing.T) {
	passwords := []string{
		"?7biBhd9d8a4$qSX",
		"RpBMT?oNXA$A$z58",
		"K#m55J$m8BJ$@1qW",
	}

	for _, pass := range passwords {
		passHash, err := HashPassword(pass)
		assert.Nil(t, err)

		passHash2, err := HashPassword(pass)
		assert.Nil(t, err)

		assert.NotEqual(t, passHash, passHash2)

		assert.Nil(t, ComparePasswordToHash(passHash, pass))
		assert.Nil(t, ComparePasswordToHash(passHash2, pass))

		assert.NotNil(t, ComparePasswordToHash(passHash, "pass"))
		assert.NotNil(t, ComparePasswordToHash(passHash2, "pass"))
	}
}
