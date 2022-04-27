module github.com/shenqianjin/rs_sample

go 1.17

require (
	github.com/qiniu/go-sdk/v7 v7.12.0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
)

require github.com/qianjin/kodo-common v0.0.0

replace github.com/qianjin/kodo-common v0.0.0 => ../../kodo-common

require github.com/qianjin/kodo-security v0.0.0
replace github.com/qianjin/kodo-security v0.0.0 => ../../kodo-security