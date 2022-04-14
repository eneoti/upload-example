package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/eneoti/upload-example/apis"
	"github.com/eneoti/upload-example/lib/logger"
	"github.com/eneoti/upload-example/services/cloudstorage"
)

func main() {
	logger := logger.NewStdLogger()

	port := 8080

	// Init the cloudstorage client
	// Exit if not connect to cloudstorage,
	cloudstorage, err := cloudstorage.GetCloudStorage("S3", logger)
	if err != nil {
		logger.Infow("Can not init cloudstorage", "error", err)
		os.Exit(0)
	}

	service, err := apis.NewService(logger, cloudstorage)
	if err != nil {
		logger.Fatalf("Failed to create new service: %s", err)
	}

	go func() {
		logger.Infow("Server is running", "port", port)
		if err = service.Start(port); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Fatalw("Failed to start the service", "error", err)
			}
		}
	}()

	// Setup the interrupt handler to gracefully exit.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGQUIT)

	go func() {
		for {
			sig := <-c

			switch sig {
			case syscall.SIGQUIT:
				buf := make([]byte, 1<<20)
				stacklen := runtime.Stack(buf, true)
				logger.Errorw("SIGQUIT", "stack", string(buf[:stacklen]))

			default:
				service.Shutdown()
				os.Exit(0)
			}
		}
	}()

	runtime.Goexit()
}
