package main

import (
	"fmt"
	"strings"
)

func main() {
	var upHosts []string
	fmt.Println(strings.Join(upHosts, ",") + "==")
	fmt.Println(strings.Join(upHosts, ",")=="")
}
