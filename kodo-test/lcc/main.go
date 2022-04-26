package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	fmt.Println(time.Now().Add(-8 * time.Hour).Unix())
	fmt.Println(time.Unix(time.Now().Add(-8 * time.Hour).Unix(), 0))
	fmt.Println("----------------")
	fmt.Println(time.Unix(1636525322, 0))
	fmt.Println(time.Unix(1636352522, 0))
	fmt.Println(time.Now().Unix())
	fmt.Println(time.Now().Unix() + 24 * 3600)
}

func TimeToUnixTimestamp(time time.Time) int64 {
	return time.Unix()
	//return fmt.Sprintf("v3:d%012d", time)
}

func UnixTimestampToTime(timestamp string) time.Time {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(i, 0)
}
