package multimod

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
)

var OuterVar1 = "OuterVar1"

func Fun1() {
	mac := auth.New("", "")
	fmt.Println(mac)
}