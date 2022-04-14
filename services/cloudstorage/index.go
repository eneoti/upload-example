package cloudstorage

import (
	"fmt"
	"upload-example/lib/logger"
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
