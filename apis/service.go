package apis

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
	"upload-example/lib/logger"
	"upload-example/services/cloudstorage"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

const (
	PAYLOAD_SIZE_LIMIT = 10240 // 10KB
)

type Service struct {
	mu                 sync.Mutex
	httpServer         http.Server
	logger             logger.Logger
	cloudstorageClient cloudstorage.ICloudStorage
}
type PayloadItem struct {
	Timestamp   time.Time   `json:"timestamp" validate:"required"`
	Name        string      `json:"name" validate:"required"`
	RequestId   string      `json:"requestId" validate:"required"`
	Context     interface{} `json:"context" validate:"required"`
	WriterKey   string      `json:"writerKey" validate:"required"`
	AnonymousId string      `json:"anonymousId" validate:"required"`
}
type Payload struct {
	Batch []PayloadItem `json:"batch" validate:"required,dive,required"`
}

func NewService(logger logger.Logger, cloudstorageClient cloudstorage.ICloudStorage) (*Service, error) {
	service := &Service{
		logger:             logger,
		cloudstorageClient: cloudstorageClient,
	}
	// Init routing
	routing := http.NewServeMux()

	// Healthz API
	routing.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Write([]byte("health check\n"))
			return
		}

		http.Error(w, "Not Found", http.StatusNotFound)
	})

	// Uploading data API
	routing.HandleFunc("/user/batch", service.uploadingData)

	service.httpServer = http.Server{
		Handler: routing,
	}

	return service, nil
}

// Start begin the http server at given port
func (s *Service) Start(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	err = s.httpServer.Serve(l)

	return err
}

// Graceful shutdown
func (s *Service) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

func (s *Service) uploadingData(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer func() {
		s.logger.Infow("Finish uploading data")
		s.mu.Unlock()
	}()

	s.logger.Infow("Uploading data is begining")
	if r.Method != http.MethodPost {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	// Limit the payload of request
	r.Body = http.MaxBytesReader(w, r.Body, PAYLOAD_SIZE_LIMIT)

	// Get the payload as JSON
	payload := Payload{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		s.logger.Errorf("Can not decode the payload:%v", err)
		http.Error(w, "Can not decode the payload", http.StatusBadRequest)
		return
	}

	// Validate the payload
	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	//TODO: Save to database

	// Upload the payload to cloud storage (S3)
	fileName := fmt.Sprintf("user-%s-%v", uuid.New().String(), time.Now().Format("2006-01-02T15:04:05"))
	buffer, _ := ioutil.ReadAll(r.Body)
	if err := s.cloudstorageClient.Upload(buffer, fileName); err != nil {
		s.logger.Infow("Uploading to cloudstorage", "error", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
}
