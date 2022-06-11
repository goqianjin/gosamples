package main

import (
	"fmt"
	"net/url"
)

func main() {
	var stringIntMap map[string]int64 /*创建集合 */
	stringIntMap = make(map[string]int64)
	fmt.Println(stringIntMap["kodo"])
	fmt.Println(stringIntMap["kodo"] + 1)
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

	siteUrl := "https://file-examples.com/storage/fef456d9a1627440e9d1c9f/2017/02/file_example_JSON_1kb.json"

	u := url.PathEscape(siteUrl)
	fmt.Println(u)
	u = url.QueryEscape(siteUrl)
	fmt.Println(u)

	strMap := map[string]string{
		"a": "1",
		"b": "2",
	}
	fmt.Println(strMap)

	updateStrMapFunc := func(m map[string]string) {
		m["C"] = "3"
	}
	updateStrMapFunc(strMap)
	fmt.Println(strMap)

}
