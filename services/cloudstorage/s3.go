package cloudstorage

import (
	"bytes"
	"fmt"
	"upload-example/lib/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	S3_REGION = "us-east-2"
	S3_BUCKET = "example"
	S3_ACL    = "public-read"
)

type S3Handler struct {
	Session *session.Session
	Bucket  string
	logger  logger.Logger
}

func NewS3(logger logger.Logger) (*S3Handler, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(S3_REGION)})
	if err != nil {
		return nil, fmt.Errorf("s3.Init, err: %v", err)
	}

	handler := S3Handler{
		logger:  logger,
		Session: sess,
		Bucket:  S3_BUCKET,
	}
	return &handler, nil

}

// TODO: This is the draft function. Need to test to make it work with S3
func (h *S3Handler) Upload(buffer []byte, fileName string) error {

	_, err := s3.New(h.Session).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(h.Bucket),
		Key:                aws.String(fileName),
		ACL:                aws.String(S3_ACL),
		Body:               bytes.NewReader(buffer),
		ContentDisposition: aws.String("Example"),
	})
	return err
}
