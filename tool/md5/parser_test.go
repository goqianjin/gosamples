package md5

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	//base64Str := "Eze7h869ksi35KJ64f4BbA=="
	base64Str := "zN9JQ31S1ozDJZMBj+buzQ=="
	bytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		fmt.Errorf("failed to decodestring. %v", err)
	}
	xstr := hex.EncodeToString(bytes)
	fmt.Println(xstr)

}
