package grpc

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/tron-us/go-common/v2/log"
)

// GracefulTerminateDetect catches SIGTERM and SIGINT
// This function only detects the signal. It does not do anything with the detect.
// Use GracefulTerminateExec func to cleanup else write your own after calling GracefulTerminateDetect
func GracefulTerminateDetect() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	select {
	case <-interrupt:
		log.Info("Received kill signal")
		break
	}
}

// GracefulTerminateExec performs a default list of tasks to clean up
// Override this function in your package to do more.
func GracefulTerminateExec(grpcServer *grpc.Server) {
	stopped := make(chan struct{})

	go func() {
		if grpcServer != nil {
			log.Info("Shutting down server gracefully")
			grpcServer.GracefulStop()
			close(stopped)
		}
	}()

	// Set timer for 10 seconds for graceful termination.
	t := time.NewTimer(10 * time.Second)
	select {
	case <-t.C:
		// Stop if graceful timer exceeds.
		log.Info("Graceful termination exceeded timeout. Stopping server!")
		grpcServer.Stop()
	case <-stopped:
		t.Stop()
	}
}
