package crypto

import (
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
