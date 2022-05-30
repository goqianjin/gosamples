package bucketcrud

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/qianjin/kodo-common/auth"
	"github.com/qianjin/kodo-common/client"
	"github.com/qianjin/kodo-sample/bucket/bucketconfig"
	"github.com/stretchr/testify/assert"
)

func TestSetDomain(t *testing.T) {
	bucketconfig.SetupEnv("192.168.200.111:10221", "192.168.200.111:10221")

	cli := client.NewManageClientWithHost("192.168.200.111:10221").
		WithKeys("6bTruscomN9sZSfq6WReU5hGDWYMDxKh0Tss4Fpr", "Q8Gn0pwsZCYEcvWewiOgNqBOCNVBs--Dqsj0cgMK").
		WithSignType(auth.SignTypeQiniu)

	bucket, resp := SetDomain(cli)
	fmt.Println("created bucket: " + bucket)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
