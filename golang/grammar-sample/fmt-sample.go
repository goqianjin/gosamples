package main

import "fmt"

func main() {
	var cmd []interface{}
	cmd = []interface{}{"Cmd1", "Cmd2"}
	data := fmt.Sprintln(cmd...)
	fmt.Println("**" + data + "**")
	fmt.Println("++" + data[:len(data)-1] + "++")

	data = fmt.Sprint(cmd...)
	fmt.Println("**" + data + "**")
	fmt.Println("++" + data[:len(data)-1] + "++")
}
