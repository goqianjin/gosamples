package rspub

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-security/kodokey"
)

func TestStat_dev(t *testing.T) {
	bucket := "kodoimport-multirs-src"
	key := "test01.txt"
	body := ""
	//var Hosts = []string{"http://10.200.20.23:9433", "http://10.200.20.23:9433"}
	path := "/stat/" + base64.URLEncoding.EncodeToString([]byte(bucket+":"+key))

	cli := client.NewClientWithHost("10.200.20.23:9433").
		WithKeys(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).
		WithSignType(auth.SignTypeQiniu)
	req := client.NewReq(http.MethodPost, path).
		RawQuery("").
		AddHeader("Host", "10.200.20.23:9433").
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		BodyStr(body).
		AddHeader("X-Qiniu-Date", strconv.FormatInt(time.Now().Add(20*time.Minute).Unix(), 10))
	resp := cli.Call(req)
	fmt.Printf("Resp: %+v\n", resp)
}
