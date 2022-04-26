package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	err := syscall.Mkdir("./run/auditlog/logupload2", 0755)
	fmt.Println(err)
	err = os.Mkdir("./run/auditlog/logupload2-1", 0755)
	fmt.Println(err)
	err = os.MkdirAll("./run/auditlog/logupload3", 0755)
	fmt.Println(err)
}
