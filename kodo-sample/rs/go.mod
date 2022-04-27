module github.com/qianjin/kodo-sample/rs

go 1.17


require github.com/qianjin/kodo-common v0.0.0

replace github.com/qianjin/kodo-common v0.0.0 => ../../kodo-common

require github.com/qianjin/kodo-security v0.0.0

replace github.com/qianjin/kodo-security v0.0.0 => ../../kodo-security

require github.com/qianjin/kodo-sample/bucket v0.0.0

replace github.com/qianjin/kodo-sample/bucket v0.0.0 => ../bucket