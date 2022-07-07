package main

import (
	"encoding/base64"
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
	fmt.Println(hex.EncodeToString([]byte("kodoimport-multiio-grammar-sample")) +
		"-" + strconv.FormatInt(int64(1380469264), 36) + ".z0.grammar-sample.qbox.me")
	fmt.Println(hex.EncodeToString([]byte("kodoimport-multiio-grammar-sample")) +
		"-" + "HwFOxpYCQU6oXoZXFOTh1mq5ZZig6Yyocgk3BTZZ" + ".z0.grammar-sample.qbox.me")

	fmt.Println("*********************")
	encodedKey := "aW9fdjI6OGdmYnBuOjE2MjE5MTA5NjE5MzgtMTYyMTkxMDk2ODI1MC50cw=="
	decodedKey, err := base64.URLEncoding.DecodeString(encodedKey)
	fmt.Println(string(decodedKey), err)
	key := "1621910961938-1621910968250.ts"
	itbl, _ := ParseItblId("8gfbpn")
	keyc := "io_v2:" + MakeItblId(itbl) + ":" + key
	encodedKey = base64.URLEncoding.EncodeToString([]byte(keyc))
	fmt.Println(encodedKey)
	// -------
	key = "errno-404"
	itbl = 480410128
	keyc = "io_v2:" + MakeItblId(itbl) + ":" + key
	fmt.Println(keyc)
	encodedKey = base64.URLEncoding.EncodeToString([]byte(keyc))
	fmt.Println(encodedKey)
	encodedKey = base64.StdEncoding.EncodeToString([]byte(keyc))
	fmt.Println(encodedKey)

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
