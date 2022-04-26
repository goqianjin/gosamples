package json

import (
	"encoding/json"
	"fmt"
)

func ToJson(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal value to json. v: %+v\n", v))
	}
	return string(bytes)
}
