package sdksample

import (
	"fmt"
	"testing"

	"github.com/qianjin/kodo-common/authkey"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

func TestBucketCreate(t *testing.T) {
	bucket := "qianjin-20220427"
	cfg := storage.Config{}
	mac := auth.New(authkey.Prod_Key_shenqianjin.AK, authkey.Prod_Key_shenqianjin.SK)
	bucketManger := storage.NewBucketManager(mac, &cfg)

	err := bucketManger.CreateBucket(bucket, "z0")
	if err != nil {
		fmt.Println(err)
	}

}
