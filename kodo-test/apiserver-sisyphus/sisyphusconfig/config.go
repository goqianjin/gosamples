package sisyphusconfig

import (
	"fmt"
	"strings"
	"time"

	"github.com/qianjin/kodo-common/env"
)

var Env env.EnvInfo

func SetupEnv(domain, host string) env.EnvInfo {
	env := env.EnvInfo{Domain: domain, Host: host}
	Env = env
	return env
}

func GenerateTaskName() string {
	RFC3339Nano := "20060102150405.999999999"
	bucketPattern := "qianjin-task-%s"
	return fmt.Sprintf(bucketPattern, strings.ReplaceAll(time.Now().Format(RFC3339Nano), ".", ""))
}
