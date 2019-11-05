package utils

import (
	"context"
	"github.com/stretchr/testify/assert"
	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
	"testing"
	"time"
)

func TestCheckRuntime(t *testing.T) {
	s := &Server{}
	ctx , cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	shared := new(sharedpb.RuntimeInfoRequest)
	runtimeInfoReport , _ := s.CheckRuntime(ctx, shared)
	assert.Equal(t, runtimeInfoReport.Status, sharedpb.RuntimeInfoReport_RUNNING , "CheckRuntime failed DB status check")
	assert.Equal(t, runtimeInfoReport.DbStatusExtra, []byte(DBConnectionHealthy), "CheckRuntime failed DB status extra check ")
}