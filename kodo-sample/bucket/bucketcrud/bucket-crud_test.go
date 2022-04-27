package bucketcrud

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
	"github.com/qianjin/kodo-security/kodokey"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	bucketconfig.SetupEnv("10.200.20.25:10221", "10.200.20.25:10221")

	cli := client.NewClientWithHost(bucketconfig.Env.Domain).
		WithKeys(kodokey.Dev_AK_general_storage_011, kodokey.Dev_SK_general_torage_011).
		WithSignType(auth.SignTypeQiniu)

	bucket, resp := Create(cli)
	fmt.Println("created bucket: " + bucket)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
