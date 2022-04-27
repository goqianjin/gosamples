package main

import (
	"encoding/json"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var json2 = jsoniter.ConfigCompatibleWithStandardLibrary

type Userinfo struct {
	user string `json: user`
	age  int    `json: age`
}

func main() {

	str := `{"user": "zhangsan", "age": 10}`
	useri := Userinfo{}
	json.Unmarshal([]byte(str), &useri)
	fmt.Println(useri)
	useri2 := Userinfo{}
	json2.UnmarshalFromString(str, &useri2)
	fmt.Println(useri2)

	Test_private_fields()
}

func Test_private_fields() {
	type TestObject struct {
		field1 string
	}
	extra.SupportPrivateFields()
	obj := TestObject{}
	fmt.Println(jsoniter.UnmarshalFromString(`{"field1":"Hello"}`, &obj))
	fmt.Println("Hello" == obj.field1)
}
