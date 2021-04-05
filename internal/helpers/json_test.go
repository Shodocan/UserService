package helpers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToJson(t *testing.T) {
	obj := map[string]interface{}{"teste": 123}
	expect, _ := json.Marshal(obj)
	result := ToJSON(obj)
	assert.Equal(t, string(expect), result, "should be an json like marshal return")
}
