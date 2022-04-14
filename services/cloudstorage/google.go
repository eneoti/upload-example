package cloudstorage

import (
	"context"
	"fmt"
	"upload-example/lib/logger"

	"cloud.google.com/go/storage"
)

const (
	Email      = "admin@upload-example.com"
	GCS_BUCKET = "example"
	Key        = "key"
)

type GCSHandler struct {
	logger        logger.Logger
	BucketHandler *storage.BucketHandle
}

func NewGCS(logger logger.Logger) (*GCSHandler, error) {
	bucketHandler, err := getGCS(GCS_BUCKET)
	if err != nil {
		return nil, fmt.Errorf("gcs.get, err: %v", err)
	}

	handler := GCSHandler{
		logger:        logger,
		BucketHandler: bucketHandler,
	}
	return &handler, nil

}

// Get Google bucket
func getGCS(bucketName string) (service *storage.BucketHandle, err error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("gcs.write file, err: %v", err)
	}
	bucket := client.Bucket(bucketName)
	return bucket, nil

}

// Upload a file to Google Cloud Storage
func (g *GCSHandler) Upload(buffer []byte, fileName string) (err error) {
	// TODO: Implement GCS uploading
	return nil
}
