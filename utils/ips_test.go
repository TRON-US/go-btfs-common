package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsLocalIp(t *testing.T) {
	isLocal, err := IsLocalIp("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, isLocal)
}
