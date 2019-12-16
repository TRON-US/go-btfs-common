package utils

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/tron-us/go-common/v2/log"
	"google.golang.org/grpc"
)

// ServerCleanupInfo contains information about grpc server
// Can be used to cleanup the buffers
type ServerCleanupInfo struct {
	GrpcServer *grpc.Server
}

// GracefulTerminateDetect catches SIGTERM and SIGINT
// This function only detects the signal. It does not do anything with the detect.
// Use GracefulTerminateExec func to cleanup else write your own after calling GracefulTerminateDetect
func (grpcInfo ServerCleanupInfo) GracefulTerminateDetect() {
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
func (grpcInfo ServerCleanupInfo) GracefulTerminateExec() {
	if grpcInfo.GrpcServer != nil {
		log.Info("Shutting down server gracefully")
		grpcInfo.GrpcServer.GracefulStop()
	}
}
