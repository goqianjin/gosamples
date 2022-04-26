package main

import "fmt"


func main() {
	zhangsan := User{"zhangsan", 31, "Beijing"}
	fmt.Println("普通输出:", zhangsan)
	fmt.Printf("v输出：%v\n", zhangsan)
	fmt.Printf("+v输出：%+v\n", zhangsan)
	fmt.Printf("#v输出：%#v\n", zhangsan)
	fmt.Printf("T输出：%T\n", zhangsan)
	fmt.Printf("b输出：%b\n", zhangsan)
	fmt.Printf("o输出：%o\n", zhangsan)
	fmt.Printf("d输出：%d\n", zhangsan)
	fmt.Printf("x输出：%x\n", zhangsan)
	fmt.Printf("X输出：%X\n", zhangsan)
	fmt.Printf("U输出：%U\n", zhangsan)
	fmt.Printf("f输出：%f\n", zhangsan)
	fmt.Printf("p输出：%p\n", zhangsan)
}
