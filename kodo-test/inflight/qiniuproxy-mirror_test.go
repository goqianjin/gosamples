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

func TestKODO14801114_KeepAuth401_UC_Image_dev(t *testing.T) {

	ioconfig.SetupEnv("127.0.0.1:7222", "127.0.0.1:7222")

	cli := client.NewClientWithHost(ioconfig.Env.Domain).
		WithKeys(authkey.Dev_Key_general_storage_011.AK, authkey.Dev_Key_general_storage_011.SK).
		WithSignType(auth.SignTypeQiniu)

	rawQuery := "bucket=qianjin-bucket-20220511143051406125&clientHash=true&key=edb2878fvodtransgzp1253922718/e7e8897d3701925919522206504/v.f230.ts&needUpload=true&rawQuery=start=0&end=2571275&type=mpegts&uid=1380538466"
	//a := `bucket=image\u0026clientHash=true\u0026key=avatar%2Fhttps%3A%2F%2Fthirdwx.qlogo.cn%2Fmmopen%2Fvi_32%2FDYAIOgq83ertk75alPq0yKoCbMnRhhPLlCkPCtNGQdPibunZLTs7ibLnTzpGpBJOLMEdg4I3EbJSdfMJ9KKu0eIg%2F132\u0026needUpload=true\u0026rawQuery=imageView2%2F1%2Fh%2F240%2Fw%2F240%2Fq%2F90\u0026uid=1380772729`
	//fmt.Println(url.QueryUnescape(url.QueryUnescape(a)))

	mirrorCofig := `{"source":"http://221.227.232.192","host":"tx-jy.snmcoocaa.aisee.tv","expires":3600,"sources":[{"addr":"http://221.227.232.192","weight":1,"backup":false}],"source_mode":2,"fragment_opt":{"fragment_size":4194304,"ignore_etag_check":false}}`
	enQuery := rawQuery //url.QueryEscape(rawQuery)
	reqQ := client.NewReq(http.MethodPost, "/mirror").
		RawQuery(enQuery).
		AddHeader("Host", ioconfig.Env.Host).
		AddHeader("Content-Type", "application/x-www-form-urlencoded").
		AddHeader("X-QN-Config", base64.URLEncoding.EncodeToString([]byte(mirrorCofig))).
		BodyStr("")
	respBo := iomodel.FetchResp{}
	resp := cli.CallWithRet(reqQ, &respBo)
	fmt.Println(resp.StatusCode)
	return

}
