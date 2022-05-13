package env

type EnvInfo struct {
	Host   string
	Domain string
}

// Prod 资源管理相关的默认域名
const (
	HostDefaultUc     = HostZ0Uc
	HostDefaultUpFast = HostZ0UpFast
	HostDefaultUp     = HostZ0Up
	HostDefaultIo     = HostZ0Io
	HostDefaultRs     = HostZ0Rs
	HostDefaultRsf    = HostZ0Rsf
	HostDefaultApi    = HostZ0Api
)

// 华东域名
const (
	HostZ0Uc     = "uc.qbox.me"        // 空间管理：http(s)://uc.qbox.me
	HostZ0UpFast = "upload.qiniup.com" // 加速上传 ：http(s)://upload.qiniup.com
	HostZ0Up     = "up.qiniup.com"     // 源站上传：http(s)://up.qiniup.com
	HostZ0Io     = "iovip.qbox.me"     // 源站下载：http(s)://iovip.qbox.me
	HostZ0Rs     = "rs.qbox.me"        // 对象管理：http(s)://rs.qbox.me
	HostZ0Rsf    = "rsf.qbox.me"       // 对象列举：http(s)://rsf.qbox.me
	HostZ0Api    = "api.qiniu.com"     // 计量查询：http(s)://api.qiniu.com
)

// 华北域名
const (
	HostZ1Uc     = HostZ0Uc               // 空间管理：http(s)://uc.qbox.me
	HostZ1UpFast = "upload-z1.qiniup.com" // 加速上传：http(s)://upload-z1.qiniup.com
	HostZ1Up     = "up-z1.qiniup.com"     // 源站上传：http(s)://up-z1.qiniup.com
	HostZ1Io     = "iovip-z1.qbox.me"     // 源站下载：http(s)://iovip-z1.qbox.me
	HostZ1Rs     = "rs-z1.qbox.me"        // 对象管理：http(s)://rs-z1.qbox.me
	HostZ1Rsf    = "rsf-z1.qbox.me"       // 对象列举：http(s)://rsf-z1.qbox.me
	HostZ1Api    = HostZ0Api              // 计量查询：http(s)://api.qiniu.com
)

// 华南域名
const (
	HostZ2Uc     = HostZ0Uc               // 空间管理：http(s)://uc.qbox.me
	HostZ2UpFast = "upload-z2.qiniup.com" // 加速上传：http(s)://upload-z2.qiniup.com
	HostZ2Up     = "up-z2.qiniup.com"     // 源站上传：http(s)://up-z2.qiniup.com
	HostZ2Io     = "iovip-z2.qbox.me"     // 源站下载：http(s)://iovip-z2.qbox.me
	HostZ2Rs     = "rs-z2.qbox.me"        // 对象管理：http(s)://rs-z2.qbox.me
	HostZ2Rsf    = "rsf-z2.qbox.me"       // 对象列举：http(s)://rsf-z2.qbox.me
	HostZ2Api    = HostZ0Api              // 计量查询：http(s)://api.qiniu.com
)

// 北美域名
const (
	HostNa0Uc     = HostZ0Uc                // 空间管理：http(s)://uc.qbox.me
	HostNa0UpFast = "upload-na0.qiniup.com" // 加速上传 ：http(s)://upload-na0.qiniup.com
	HostNa0Up     = "up-na0.qiniup.com"     // 源站上传：http(s)://up-na0.qiniup.com
	HostNa0Io     = "iovip-na0.qbox.me"     // 源站下载：http(s)://iovip-na0.qbox.me
	HostNa0Rs     = "rs-na0.qbox.me"        // 对象管理：http(s)://rs-na0.qbox.me
	HostNa0Rsf    = "rsf-na0.qbox.me"       // 对象列举：http(s)://rsf-na0.qbox.me
	HostNa0Api    = HostZ0Api               // 计量查询：http(s)://api.qiniu.com
)

// 东南亚域名
const (
	HostAs0Uc     = HostZ0Uc                // 空间管理：http(s)://uc.qbox.me
	HostAs0UpFast = "upload-as0.qiniup.com" // 加速上传：http(s)://upload-as0.qiniup.com
	HostAs0Up     = "up-as0.qiniup.com"     // 源站上传：http(s)://up-as0.qiniup.com
	HostAs0Io     = "iovip-as0.qbox.me"     // 源站下载：http(s)://iovip-as0.qbox.me
	HostAs0Rs     = "rs-as0.qbox.me"        // 对象管理：http(s)://rs-as0.qbox.me
	HostAs0Rsf    = "rsf-as0.qbox.me"       // 对象列举：http(s)://rsf-as0.qbox.me
	HostAs0Api    = HostZ0Api               // 计量查询：http(s)://api.qiniu.com
)

// 华东-浙江2域名
const (
	HostCnEast2Uc     = HostZ0Uc                      // 空间管理：http(s)://uc.qbox.me
	HostCnEast2UpFast = "upload-cn-east-2.qiniup.com" // 加速上传：http(s)://upload-cn-east-2.qiniup.com
	HostCnEast2Up     = "up-cn-east-2.qiniup.com"     // 源站上传：http(s)://up-cn-east-2.qiniup.com
	HostCnEast2Io     = "iovip-cn-east-2.qiniuio.com" // 源站下载：http(s)://iovip-cn-east-2.qiniuio.com
	HostCnEast2Rs     = "rs-cn-east-2.qiniuapi.com"   // 对象管理：http(s)://rs-cn-east-2.qiniuapi.com
	HostCnEast2Rsf    = "rsf-cn-east-2.qiniuapi.com"  // 对象列举：http(s)://rsf-cn-east-2.qiniuapi.com
	HostCnEast2Api    = HostZ0Api                     // 计量查询：http(s)://api.qiniu.com
)
