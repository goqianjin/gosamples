package main

import (
	"flag"
	"fmt"
)

func init() {
	fmt.Println("init 1")
}

type User struct {
	Name string
	Age int8
	Address string
}

func main() {
	fmt.Println(3/2)
	zhangsan := User{"zhangsan", 10, "Shanghai"}
	fmt.Printf("%+v\n", zhangsan)
	updateUser(zhangsan)
	fmt.Printf("%+v\n", zhangsan)
	updateUserByPointer(&zhangsan)
	fmt.Printf("%+v\n", zhangsan)

	// 测试输出格式
	fmt.Printf("%p\n", &zhangsan)

	confName := flag.String("f", "kodoimport.conf", "conf file")

	fmt.Println("Use the config file of ", *confName)
	fmt.Println(confName)

}

// update by pure param
func updateUser(User User) {
	User.Age += 5
}
// update by pure param
//func updateUser(User User, Index int) {
//	User.Age += 5
//}

// update by pointer ref
func updateUserByPointer(User *User) {
	User.Age += 5
}


