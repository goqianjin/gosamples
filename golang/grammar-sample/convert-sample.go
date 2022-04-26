package main

import (
	"fmt"
	"strconv"
)

func main() {
	var intValue int = 1
	//fmt.Println(intValue)
	fmt.Println("整型-整型转换..........")
	// int to int32
	var int32Value int32 = int32(intValue)
	fmt.Println(int32Value)
	// int32 to int
	intValue = int(int32Value)
	fmt.Println(intValue)
	// int to int64
	var int64Value int64 = int64(intValue)
	fmt.Println(int64Value)
	// int64 to int
	intValue = int(int64Value)
	fmt.Println(intValue)
	// int32 to int64
	int64Value = int64(int32Value)
	fmt.Println(int64Value)
	// int64 to int32
	int32Value = int32(int64Value)
	fmt.Println(int32Value)

	fmt.Println("浮点型-浮点型转换..........")
	var float32Value float32 = 1.2
	// float32 to float64
	var float64Value float64 = float64(float32Value)
	fmt.Println(float64Value)
	// float64 to float32
	float32Value = float32(float64Value)
	fmt.Println(float32Value)

	fmt.Println("整型-浮点型转换..........")
	float64Value = 123.456
	// float64 to int64
	int64Value = int64(float64Value)
	fmt.Println(int64Value)
	// float64 to int32
	int32Value = int32(float64Value)
	fmt.Println(int32Value)
	// float64 to int
	intValue = int(float64Value)
	fmt.Println(intValue)
	//
	float32Value = 111.222
	// float32 to int64
	int64Value = int64(float32Value)
	fmt.Println(int64Value)
	// float32 to int32
	int32Value = int32(float32Value)
	fmt.Println(int32Value)
	// float32 to int
	intValue = int(float32Value)
	fmt.Println(intValue)
	// int to float64
	float64Value = float64(intValue)
	fmt.Println(float64Value)
	// int32 to float64
	float64Value = float64(int32Value)
	fmt.Println(float64Value)
	// int64 to float64
	float64Value = float64(int64Value)
	fmt.Println(float64Value)
	// int to float32
	float32Value = float32(intValue)
	fmt.Println(float32Value)
	// int32 to float32
	float32Value = float32(int32Value)
	fmt.Println(float32Value)
	// int64 to float32
	float32Value = float32(int64Value)
	fmt.Println(float32Value)

	// int string
	var stringVar string
	fmt.Println("整型-字符串转换..........")
	// int to string
	stringVar = strconv.Itoa(intValue)
	fmt.Println(stringVar)
	// method 2 -- TODO dig into
	stringVar = strconv.FormatInt(int64(intValue), 10)
	fmt.Println(stringVar)
	// int32 to string
	stringVar = strconv.FormatInt(int64(int32Value), 10)
	fmt.Println(stringVar)
	// int64 to string
	stringVar = strconv.FormatInt(int64Value, 10)
	fmt.Println(stringVar)
	// string to int
	intValue, _ = strconv.Atoi(stringVar)
	fmt.Println(intValue)
	// method 2
	int64Value, _ = strconv.ParseInt(stringVar, 10, 32)
	intValue = int(int64Value)
	fmt.Println(intValue)
	// string to int32
	int64Value, _ = strconv.ParseInt(stringVar, 10, 64)
	int32Value = int32(int64Value)
	fmt.Println(intValue)
	// string to int64
	int64Value, _ = strconv.ParseInt(stringVar, 10, 64)
	fmt.Println(int64Value)


	fmt.Println("浮点型-字符串转换..........")
	// float32 to string
	stringVar = strconv.FormatFloat(float64(float32Value), 'E', -1, 32)
	fmt.Println(stringVar)
	// float64 to string
	stringVar = strconv.FormatFloat(float64Value, 'E', -1, 64)
	fmt.Println(stringVar)
	// string to float32
	float64Value, _ = strconv.ParseFloat(stringVar, 32)
	float32Value = float32(float64Value)
	fmt.Println(float32Value)
	// string to float64
	float64Value, _ = strconv.ParseFloat(stringVar, 32)
	fmt.Println(float64Value)
}
