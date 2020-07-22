package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTronAddress(t *testing.T) {
	keys, err := FromPrivateKey("CAISIJqwbU3ceD6u2tsrGWe/Zk+MhZ9mNcQWxnlB03Zv3uay")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "TEmH6bFE8FahLDTJwZqfAVuB7CaxjvPmZV", keys.Base58Address)
	assert.Equal(t, "4134970e58155b49379b5a1fdd9f06949ef6bbdf58", keys.HexAddress)
	assert.Equal(t, "9ab06d4ddc783eaedadb2b1967bf664f8c859f6635c416c67941d3766fdee6b2", keys.HexPrivateKey)
	assert.Equal(t, "CAISIQPUMOOBFVu1VVAgSF7VKUzgj66yWJmgtZs0GY3r9ZCxpw==", keys.Base64PubKey)
	assert.Equal(t, "04d430e381155bb5555020485ed5294ce08faeb25899a0b59b34198debf590b1a7ce7743714e4c8fe23b5c5a7dfe42a6163b290b5432129238a8fd1faf4a08ffcf", keys.HexPubKey)
}
