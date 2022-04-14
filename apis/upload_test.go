package apis

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/eneoti/upload-example/lib/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock the Upload function to the cloud storage.
type cloudStorageClientMock struct {
	mock.Mock
}

func (c *cloudStorageClientMock) Upload(buffer []byte, fileName string) error {
	args := c.Called(buffer, fileName)
	return args.Error(0)
}

func setup(t *testing.T) (*cloudStorageClientMock, *Service, *httptest.ResponseRecorder) {
	cloudStorageClient := new(cloudStorageClientMock)

	log := logger.NewStdLogger()

	service, err := NewService(log, cloudStorageClient)
	assert.Nil(t, err)

	w := httptest.NewRecorder()

	return cloudStorageClient, service, w
}

var anyBytes = mock.MatchedBy(func(buffer []byte) bool {
	return true
})

// Test the uploading API to be sucessfully when the cloud storage has problem.
func TestPostUploadingData_Sucess_EvenStorageCloudProblem(t *testing.T) {
	cloudStorageClient, service, w := setup(t)

	data, err := os.Open("./data-test/payload.json")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	req, err := http.NewRequest(http.MethodPost, "/user/batch", data)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	cloudStorageClient.On("Upload", anyBytes, mock.Anything).Return(
		errors.New("uploading error"),
	)

	service.httpServer.Handler.ServeHTTP(w, req)

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 200, w.Code)

	cloudStorageClient.AssertExpectations(t)
}

// Test the uploading API to be failed when the payload is invalid format.
func TestPostUpoadingData_Failed_Due2InvalidBody(t *testing.T) {
	cloudStorageClient, service, w := setup(t)

	data, err := os.Open("./data-test/invalid-payload.json")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	req, err := http.NewRequest(http.MethodPost, "/user/batch", data)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	cloudStorageClient.AssertNotCalled(t, "Upload", anyBytes, mock.Anything)

	service.httpServer.Handler.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "Invalid payload\n", w.Body.String())

}

// Test the uploading API to be failed when the payload is exceeded 10KB.
func TestPostUpoadingData_Failed_Due2Exceed10KB(t *testing.T) {
	cloudStorageClient, service, w := setup(t)
	data, err := os.Open("./data-test/payload-exceed-10KB.json")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	req, err := http.NewRequest(http.MethodPost, "/user/batch", data)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	cloudStorageClient.AssertNotCalled(t, "Upload", anyBytes, mock.Anything)
	service.httpServer.Handler.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "Can not decode the payload\n", w.Body.String())

}

// Test the uploading API to be sucessfully when the payload is less than 10KB.
func TestPostUpoadingData_Success(t *testing.T) {
	cloudStorageClient, service, w := setup(t)

	data, err := os.Open("./data-test/payload.json")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	req, err := http.NewRequest(http.MethodPost, "/user/batch", data)
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	cloudStorageClient.On("Upload", anyBytes, mock.Anything).Return(nil)

	service.httpServer.Handler.ServeHTTP(w, req)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 200, w.Code)

	cloudStorageClient.AssertExpectations(t)

}
