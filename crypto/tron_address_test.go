package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tron-us/go-common/v2/log"

	"go.uber.org/zap"
)

func TestChangeFormat(t *testing.T) {
	peerId := "16Uiu2HAkxrYc8rkZFWF2RfAW2jdu1zYpruEJNkjTk6nBsVhzx2x4"
	pubAddress, err := GetTronPubKeyFromPeerIdPretty(peerId)
	assert.Equal(t, err, nil, "err == nil")
	assert.Equal(t, *pubAddress, "TK33CBqWDuQzhkDJiJRk8BpyKDvVPqxPrZ", "*pubAddress == TK33CBqWDuQzhkDJiJRk8BpyKDvVPqxPrZ")
}

func TestGetAddress(t *testing.T) {
	peerIds := []string{
		"16Uiu2HAkzhsnMffxirhJ39ZMdPcXdrcbfrZsvcQdUeVwSBzBqBK7",
		"16Uiu2HAm1yEfFmzC1enfBcfbwf51YA15e4tRd9VRT65TELA1ykAD",
		"16Uiu2HAmCAqAfWHnpYKqvx627K1C8Bm3ixJiwALXwXWGpRBLH8dz",
	}

	for _, id := range peerIds {
		pubAddress, err := GetTronPubKeyFromPeerIdPretty(id)
		assert.Equal(t, err, nil, "err == nil")

		log.Info("get", zap.String("peerid", id), zap.String("address", *pubAddress))
	}

}
