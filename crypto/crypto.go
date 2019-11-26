package crypto

import (
	"encoding/base64"

	"github.com/gogo/protobuf/proto"
	ic "github.com/libp2p/go-libp2p-core/crypto"
)

func Sign(key ic.PrivKey, channelMessage proto.Message) ([]byte, error) {
	raw, err := proto.Marshal(channelMessage)
	if err != nil {
		return nil, err
	}
	return key.Sign(raw)
}

func Verify(key ic.PubKey, channelMessage proto.Message, sig []byte) (bool, error) {
	raw, err := proto.Marshal(channelMessage)
	if err != nil {
		return false, err
	}
	return key.Verify(raw, sig)
}

// private key string to ic.PrivKey interface
// btfs config stores base64 of private key
func ToPrivKey(privKey string) (ic.PrivKey, error) {
	raw, err := base64.StdEncoding.DecodeString(privKey)
	if err != nil {
		return nil, err
	}
	return ic.UnmarshalPrivateKey(raw)
}

// public key string to ic.PubKey interface
func ToPubKey(pubKey []byte) (ic.PubKey, error) {
	return ic.UnmarshalSecp256k1PublicKey(pubKey)
}
