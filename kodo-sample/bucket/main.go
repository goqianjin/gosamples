package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
)

// username: general_storage_002@test.qiniu.io, "uid" : 1380469264, "itbl" : 514872649

func main() {
	fmt.Println(MakeItblId(514872649))
	fmt.Println(ParseItblId(MakeItblId(514872649)))
	fmt.Println(ParseItblId("8ijqyq"))
	fmt.Println(MakeItblId(514883906))
	fmt.Println(ParseItblId("8chgd0"))
	//
	fmt.Println(hex.EncodeToString([]byte("kodoimport-multiio-grammar-sample"))+
		"-" + strconv.FormatInt(int64(1380469264), 36) + ".z0.grammar-sample.qbox.me")
	fmt.Println(hex.EncodeToString([]byte("kodoimport-multiio-grammar-sample"))+
		"-" + "HwFOxpYCQU6oXoZXFOTh1mq5ZZig6Yyocgk3BTZZ" + ".z0.grammar-sample.qbox.me")


	fmt.Println("*********************")
	fmt.Println(MakeItblId(483647306))

}


func MakeItblId(itbl uint32) string {
	return strconv.FormatInt(int64(itbl), 36)
}

// ------------------------------------------------------------------------

func ParseItblId(id string) (itbl uint32, err error) {
	v, err := strconv.ParseUint(id, 36, 32)
	if err != nil {
		return 0, errors.New("Invalid itbl: " + id)
	}
	return uint32(v), nil
}