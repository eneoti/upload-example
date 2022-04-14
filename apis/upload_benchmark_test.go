package apis

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eneoti/upload-example/lib/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupBenchmark(b *testing.B) (*cloudStorageClientMock, *Service, *httptest.ResponseRecorder) {
	cloudStorageClient := new(cloudStorageClientMock)

	log := logger.NewStdLogger()

	service, err := NewService(log, cloudStorageClient)
	assert.Nil(b, err)

	w := httptest.NewRecorder()

	return cloudStorageClient, service, w
}

// Benchmark the uploading data API
func BenchmarkUploadingAPITest(b *testing.B) {
	cloudStorageClient, service, w := setupBenchmark(b)
	// assume that the cloudstorage take 01 seconds to hanlde the request.
	cloudStorageClient.On("Upload", anyBytes, mock.Anything).After(1 * time.Second).Return(nil)
	// Make same payload for every request
	payload, err := ioutil.ReadFile("./data-test/payload.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	b.RunParallel(func(pb *testing.PB) {
		b.ReportAllocs()
		b.ResetTimer()
		for pb.Next() {
			req, err := http.NewRequest(http.MethodPost, "/user/batch", bytes.NewReader(payload))
			assert.Nil(b, err)
			service.httpServer.Handler.ServeHTTP(w, req)
		}
	})
}
