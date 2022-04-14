package cloudstorage

import (
	"fmt"
	"upload-example/lib/logger"

	"go.uber.org/ratelimit"
)

type ICloudStorage interface {
	// Upload base on the filePath
	Upload(buffer []byte, fileName string) error
}

// Factory the client to work with S3/GCS
func GetCloudStorage(storageType string, logger logger.Logger) (ICloudStorage, error) {
	if storageType == "S3" {
		return NewS3(logger)
	}
	if storageType == "GCS" {
		return NewGCS(logger)
	}
	return nil, fmt.Errorf("wrong storage type passed")
}

type CSEngine struct {
	cloudstorageClient interface{}
	logger             logger.Logger
	s3RateLimit        ratelimit.Limiter
	gCSRateLimit       ratelimit.Limiter
}

func NewCSEngine(logger logger.Logger, cloudstorageClient interface{}) (*CSEngine, error) {
	s3RateLimit := ratelimit.New(10)
	gCSRateLimit := ratelimit.New(10)
	csEngine := CSEngine{
		s3RateLimit:        s3RateLimit,
		gCSRateLimit:       gCSRateLimit,
		logger:             logger,
		cloudstorageClient: cloudstorageClient,
	}

	return &csEngine, nil
}
func (w *CSEngine) GetRateLimit() ratelimit.Limiter {
	switch w.cloudstorageClient.(type) {
	default:
		return w.s3RateLimit
	case S3Handler:
		return w.s3RateLimit
	case GCSHandler:
		return w.gCSRateLimit

	}
}
func (w *CSEngine) GetGCSType() string {
	switch w.cloudstorageClient.(type) {
	default:
		return "S3"
	case S3Handler:
		return "S3"
	case GCSHandler:
		return "GCS"

	}
}

func (w *CSEngine) Do(buffer []byte, fileName string) {
	rl := w.GetRateLimit()
	rl.Take()
	go func() {
		err := w.cloudstorageClient.(ICloudStorage).Upload(buffer, fileName)
		if err == nil {
			w.logger.Infof("Upload %v to %v sucessfully ", fileName, w.GetGCSType())
		}

	}()
}
