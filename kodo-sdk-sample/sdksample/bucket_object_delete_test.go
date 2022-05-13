package sdksample

import (
	"fmt"
	"testing"

	"github.com/qiniu/go-sdk/v7/client"

	"github.com/qianjin/kodo-common/authkey"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

func TestBucketObjectDelete(t *testing.T) {
	bucket := "qj-kodoimport-202203-src"
	key := "test02.txt"
	cfg := storage.Config{}
	client.DebugMode = true
	client.DeepDebugInfo = true

	mac := auth.New(authkey.Prod_Key_kodolog.AK, authkey.Prod_Key_kodolog.SK)
	bucketManger := storage.NewBucketManager(mac, &cfg)

	err := bucketManger.Delete(bucket, key)
	if err != nil {
		fmt.Println(err)
	}

}
