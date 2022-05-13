package inflight

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/qianjin/kodo-common/authkey"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/io/iomodel"

	"github.com/qianjin/kodo-sample/io/ioconfig"
)

func Test_Qiniuproxy_Mirror_Local(t *testing.T) {
	ioconfig.SetupEnv("127.0.0.1:7222", "127.0.0.1:7222")
	test_Qiniuproxy_Mirror(t, authkey.Dev_Key_general_storage_011)
}

func Test_Qiniuproxy_Mirror_Dev(t *testing.T) {
}

func Test_Qiniuproxy_Mirror_Prod(t *testing.T) {
}

func test_Qiniuproxy_Mirror(t *testing.T, authKey authkey.AuthKey) {

	cli := client.NewClientWithHost(ioconfig.Env.Domain).
		WithKeys(authKey.AK, authKey.SK).
		WithSignType(auth.SignTypeQiniu)

	rawQuery := "bucket=qianjin-bucket-20220511143051406125&clientHash=true&key=edb2878fvodtransgzp1253922718/e7e8897d3701925919522206504/v.f230.ts&needUpload=true&rawQuery=start=0&end=2571275&type=mpegts&uid=1380538466"
	mirrorCofig := `{"source":"http://221.227.232.192","host":"tx-jy.snmcoocaa.aisee.tv","expires":3600,"sources":[{"addr":"http://221.227.232.192","weight":1,"backup":false}],"source_mode":2,"fragment_opt":{"fragment_size":4194304,"ignore_etag_check":false}}`
	reqQ := client.NewReq(http.MethodPost, "/mirror").
		RawQuery(rawQuery).
		AddHeader("Host", ioconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		AddHeader("X-QN-Config", base64.URLEncoding.EncodeToString([]byte(mirrorCofig))).
		BodyStr("")
	respBo := iomodel.FetchResp{}

	resp := cli.CallWithRet(reqQ, &respBo)
	fmt.Println(resp.StatusCode)
	return
}
