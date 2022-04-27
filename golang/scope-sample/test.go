package main

import (
	"fmt"
	"scope/model1"
	"scope/model2"
)

var OutVariable1 = "out variable1"

func main() {
	fmt.Println(model1.Variable1)
	model1.Method1()
	fmt.Println(model2.Variable1)


}

func ref() {
	model1.Method1()
}
