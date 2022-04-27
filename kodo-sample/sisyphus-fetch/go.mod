module fetch

go 1.17

require github.com/pborman/uuid v1.2.1

require (
	github.com/google/uuid v1.0.0 // indirect
	github.com/qiniu/go-sdk/v7 v7.12.1 // indirect
	go.mongodb.org/mongo-driver v1.7.4 // indirect
)

require github.com/qianjin/kodo-security v0.0.0

replace github.com/qianjin/kodo-security v0.0.0 => ../../kodo-security
