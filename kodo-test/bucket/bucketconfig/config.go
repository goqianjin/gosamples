package bucketconfig

import (
	"fmt"
	"strings"
	"time"
)

var Env EnvInfo

type EnvInfo struct {
	Host   string
	Domain string
}

func SetupEnv(domain, host string) EnvInfo {
	env := EnvInfo{Domain: domain, Host: host}
	Env = env
	return env
}

func GenerateBucketName() string {
	RFC3339Nano := "20060102150405.999999999"
	bucketPattern := "qianjin-bucket-%s"
	return fmt.Sprintf(bucketPattern, strings.ReplaceAll(time.Now().Format(RFC3339Nano), ".", ""))
}
