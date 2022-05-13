module github.com/qianjin/kodo-sample/up

go 1.17

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/qiniu/go-sdk/v7 v7.12.1
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

require (
	github.com/qianjin/kodo-common v0.0.0
	github.com/stretchr/testify v1.7.1
)

replace github.com/qianjin/kodo-common v0.0.0 => ../../kodo-common

require github.com/qianjin/kodo-security v0.0.0
replace github.com/qianjin/kodo-security v0.0.0 => ../../kodo-security

require github.com/qianjin/kodo-sample/bucket v0.0.0

replace github.com/qianjin/kodo-sample/bucket v0.0.0 => ../bucket
