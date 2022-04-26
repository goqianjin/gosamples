package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)


func init() {
	fmt.Println("init 2")
}

func main() {
	RegionUpHostMap := map[string]string{
		"z0": "http://xxx.qiniu.com",
		"z1": "xxx2.qiniu.com,xxx3.qiniu.com___`,",
		"z2": "xxx2.qiniu.com,xxx3.qiniu.com",
	}
	// validate
	for zone := range RegionUpHostMap {
		upHosts := strings.Trim(RegionUpHostMap[zone], ",")
		upHosts = strings.Trim(upHosts, " ")
		upHostSlice := strings.Split(upHosts, ",")
		for index, host := range upHostSlice {
			if !strings.Contains(host, "://") {
				host = "http://" + host
			}
			// parse to validate url
			if _, err := url.Parse(host); err != nil {
				log.Fatalf("region_up_host_map field contains invalid host on region: %v", zone)
			}
			// update
			upHostSlice[index] = host
		}
		RegionUpHostMap[zone] = strings.Join(upHostSlice, ",")
		fmt.Println(RegionUpHostMap)
		sl := [] string {"aa", "bb", "cc"}
		fmt.Println(strings.Join(sl, ","))
	}
}


//func updateUserByPointer(User *User) {
//	User.Age += 5
//}