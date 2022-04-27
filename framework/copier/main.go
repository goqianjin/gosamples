package main

import (
	"fmt"

	"github.com/jinzhu/copier"
)

func main() {

	u1 := User{
		Name:        "zhangsan",
		Role:        "student",
		Age:         14,
		EmployeCode: 123456,
		Salary:      12,
		Addr: Address{
			City:    "shanghai",
			Country: "china",
			Town:    "zhangjiang",
		},
	}
	u2 := User{}
	copier.Copy(&u2, &u1)
	fmt.Printf("u1: %p --> %+v\n", &u1, u1)
	fmt.Printf("u2: %p --> %+v\n", &u2, u2)
	fmt.Printf("adress of u1: %p --> %+v\n", &u1.Addr, u1.Addr)
	fmt.Printf("adress of u2: %p --> %+v\n", &u2.Addr, u2.Addr)

}

type User struct {
	Name        string
	Role        string
	Age         int32
	EmployeCode int64 `copier:"EmployeNum"` // specify field name

	// Explicitly ignored in the destination struct.
	Salary int
	Addr   Address
}

type Address struct {
	City    string
	Country string
	Town    string
}
