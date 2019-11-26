package crypto

import (
	"testing"

	ledgerPb "github.com/tron-us/go-btfs-common/protos/ledger"
)

const (
	KeyString = "CAISIJFNZZd5ZSvi9OlJP/mz/vvUobvlrr2//QN4DzX/EShP"
)

func TestSignVerify(t *testing.T) {
	// test get privKey and pubKey
	privKey, err := ToPrivKey(KeyString)
	if err != nil {
		t.Error("ToPrivKey failed")
		return
	}

	rawPubKey, err := privKey.GetPublic().Raw()
	if err != nil {
		t.Error("get raw public key from privKey failed")
		return
	}
	pubKey, err := ToPubKey(rawPubKey)
	if err != nil {
		t.Error("ToPubKey failed")
		return
	}

	// test sign and verify the key string
	message := &ledgerPb.PublicKey{
		Key: rawPubKey,
	}

	sign, err := Sign(privKey, message)
	if err != nil {
		t.Error("Sign with private key failed")
		return
	}
	ret, err := Verify(pubKey, message, sign)
	if err != nil || !ret {
		t.Error("Verify with public key failed")
	}
}
