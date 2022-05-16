package client_s3

import "github.com/aws/aws-sdk-go/service/s3"

func main() {
	svc := s3.New(nil)
	svc.PutObject()
	svc.DeleteBucketTagging()
	svc.DeleteBucketTaggingWithContext()
}
