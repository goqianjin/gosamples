package main

import "fmt"

func main() {
	var stringIntMap map[string]int64 /*创建集合 */
	stringIntMap = make(map[string]int64)
	fmt.Println(stringIntMap["kodo"])
	fmt.Println(stringIntMap["kodo"]+1)
	stringIntMap["s1"] = 10
	stringIntMap["s2"] = 20
	stringIntMap["s3"] = 30
	fmt.Println(stringIntMap[""])
	fmt.Println(stringIntMap)
	delete(stringIntMap, "s2")
	fmt.Println(stringIntMap)
	// delete一个不存在的key
	delete(stringIntMap, "xxxx")
	fmt.Println(stringIntMap)

}
