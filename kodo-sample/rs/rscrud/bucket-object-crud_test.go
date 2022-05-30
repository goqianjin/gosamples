package rscrud

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/qianjin/kodo-sample/rs/rsmodel"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/rs/rsconfig"
	"github.com/qianjin/kodo-security/kodokey"
	"github.com/stretchr/testify/assert"
)

func TestDeleteObject(t *testing.T) {
	rsconfig.SetupEnv("10.200.20.23:9433", "10.200.20.23:9433")

	cli := client.NewManageClientWithHost(rsconfig.Env.Domain).
		WithKeys(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).
		WithSignType(auth.SignTypeQiniu)

	bucket := "qj-kodoimport-202203-src"
	key := "test01.txt"
	_, resp := DeleteObject(cli, rsmodel.DeleteObjectReq{Bucket: bucket, Key: key})
	fmt.Printf("delete bucket: %s, key: %s\n", bucket, key)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
