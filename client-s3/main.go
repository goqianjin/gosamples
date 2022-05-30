package client_s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	svc := s3.New(nil)
	svc.PutObject(nil)
	svc.DeleteBucketTagging(nil)
	svc.DeleteBucketTaggingWithContext(aws.Context{}, nil)
}
