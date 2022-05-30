package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func main() {
	type User struct {
		Id uint64 `json:id`
	}
	u := User{
		Id: uint64(0x0001),
	}
	fmt.Println(u)
	us, err := json.Marshal(u)
	fmt.Println(string(us), err)

	//uxs := `{"Id":0x0001}`
	uxs := `{"Id":321}`
	var ux User
	err = json.Unmarshal([]byte(uxs), &ux)
	fmt.Println(ux, err)
	var ud interface{}
	err = json.Unmarshal([]byte(uxs), &ud)
	fmt.Println(ud, err)
	fmt.Println("aaa ---- " + fmt.Sprint(ud))
	var i int
	is := `12`
	err = json.Unmarshal([]byte(is), &i)
	fmt.Println(i, err)
	var s string
	ss := `"123"`
	err = json.Unmarshal([]byte(ss), &s)
	fmt.Println("----->", s, err)
	ss = `""`
	err = json.Unmarshal([]byte(ss), &s)
	fmt.Println("----->", s, err)

	bis := strconv.FormatUint(uint64(0x000011), 2)
	fmt.Println(bis)
	bi, err := strconv.ParseUint(bis, 2, 64)
	fmt.Println(bi, err)
	bi, err = strconv.ParseUint("000000000001", 2, 64)
	fmt.Println(bi, err)
	bi, err = strconv.ParseUint("", 2, 64)
	fmt.Println(bi, err)
	bi, err = strconv.ParseUint("1000000000000000000000000000000000000000000000000000000000000000", 2, 64)
	fmt.Println(bi, err)
	bi, err = strconv.ParseUint("10000000000000000000000000000000000000000000000000000000000000000", 2, 64)
	fmt.Println(bi, err)
	bi, err = strconv.ParseUint("0b000000000001", 0, 64)
	fmt.Println(bi, err)

}
