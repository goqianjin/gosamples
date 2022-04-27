package env

type EnvInfo struct {
	Host   string
	Domain string
}

// PROD资源管理相关的默认域名
const (
	DefaultRsHost  = "rs.qiniu.com"
	DefaultRsfHost = "rsf.qiniu.com"
	DefaultAPIHost = "api.qiniu.com"
	DefaultPubHost = "pu.qbox.me:10200"

	DefaultUcHost = "uc.qbox.me"
)
