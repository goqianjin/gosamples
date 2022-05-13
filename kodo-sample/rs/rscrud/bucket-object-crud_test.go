package bucketcrud

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/qianjin/kodo-sample/bucket/bucketmodel"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
	"github.com/qianjin/kodo-security/kodokey"
	"github.com/stretchr/testify/assert"
)

func TestDeleteObject(t *testing.T) {
	bucketconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")

	cli := client.NewClientWithHost(bucketconfig.Env.Domain).
		WithKeys(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).
		WithSignType(auth.SignTypeQiniu)

	bucket := "qj-kodoimport-202203-src"
	key := "test01.txt"
	_, resp := DeleteObject(cli, bucketmodel.DeleteObjectReq{Bucket: bucket, Key: key})
	fmt.Printf("delete bucket: %s, key: %s\n", bucket, key)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
