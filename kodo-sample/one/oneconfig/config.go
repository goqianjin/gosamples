package oneconfig

import (
	"github.com/qianjin/kodo-common/env"
)

var Env env.EnvInfo

func SetupEnv(domain, host string) env.EnvInfo {
	env := env.EnvInfo{Domain: domain, Host: host}
	Env = env
	return env
}
