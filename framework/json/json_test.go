package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Man struct {
	Name    string
	Age     int
	Address string
}

type Student struct {
	Man
	No int
}

func TestAbc(t *testing.T) {
	s := Student{No: 2223, Man: Man{Name: "zhangsan", Age: 12}}
	fmt.Printf("ooo: %v\n", s)
	sb, _ := json.Marshal(s)
	fmt.Printf("json: %v\n", string(sb))

	s2 := Student{}
	json.Unmarshal([]byte("{\"Name\":\"zhangsan\",\"Age\":12}"), &s2)
	fmt.Printf("s2 - ooo: %v\n", s2)

	s2b, _ := json.Marshal(s2)
	fmt.Printf("json: %v\n", string(s2b))
}
