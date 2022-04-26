package main

import (
	"fmt"
	"sync"
)

func main() {
	var syncMap sync.Map
	syncMap.Store("A", 1)
	syncMap.LoadOrStore("B", 2)
	syncMap.Store("C", 3)
	fmt.Printf("%+v\n", syncMap)
	var keys []string
	syncMap.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	fmt.Println(keys)
	var maps = make(map[string]int)
	syncMap.Range(func(key, value interface{}) bool {
		maps[key.(string)] = value.(int)
		return true
	})
	fmt.Println(maps)

}
