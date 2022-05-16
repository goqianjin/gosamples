module sdk

go 1.17

require github.com/aliyun/aliyun-oss-go-sdk v2.2.2+incompatible

require (
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	golang.org/x/time v0.0.0-20220411224347-583f2d630306 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

require github.com/qianjin/kodo-security v0.0.0

replace github.com/qianjin/kodo-security v0.0.0 => ../../kodo-security
