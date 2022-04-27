package model2

import (
	"fmt"
	"scope/model1"
)

var variable1 = 22
var Variable1 = 222

func method1() {
	fmt.Println("model2 --> mehtod1")
}

func Method1() {
	fmt.Println("model2 --> Mehtod1")
}

func test() {
	model1.Method1()
	fmt.Println(model1.Variable1)
}