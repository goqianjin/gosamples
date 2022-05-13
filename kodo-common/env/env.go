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

const (
	Host_Bucket              = "uc.qbox.me"
	Host_IO                  = "iovip.qbox.me"
	Host_Huangdong_Up_Origin = "up.qiniup.com"
)
