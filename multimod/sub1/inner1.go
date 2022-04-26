package sub1

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
)

func SubFun1() {
	mac := auth.New("", "")
	fmt.Println(mac)
}

func SubFun2() {

}
