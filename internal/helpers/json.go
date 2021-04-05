package helpers

import "encoding/json"

func ToJSON(obj interface{}) string {
	result, _ := json.Marshal(obj)
	return string(result)
}
