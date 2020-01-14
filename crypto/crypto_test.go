package crypto

import (
	"bytes"
	"testing"

	ledgerPb "github.com/tron-us/go-btfs-common/protos/ledger"
)

const (
	KeyString  = "CAISIJFNZZd5ZSvi9OlJP/mz/vvUobvlrr2//QN4DzX/EShP"
	EncryptKey = "Tron2theMoon1234"
)

func TestSignVerify(t *testing.T) {
	// test get privKey and pubKey
	privKey, err := ToPrivKey(KeyString)
	if err != nil {
		t.Fatal("ToPrivKey failed", err)
	}

	rawPubKey, err := privKey.GetPublic().Raw()
	if err != nil {
		t.Fatal("Get raw public key from privKey failed", err)
	}
	pubKey, err := ToPubKeyRaw(rawPubKey)
	if err != nil {
		t.Fatal("ToPubKeyRaw failed", err)
	}

	// test sign and verify the key string
	message := &ledgerPb.PublicKey{
		Key: rawPubKey,
	}

	sign, err := Sign(privKey, message)
	if err != nil {
		t.Fatal("Sign with private key failed", err)
	}
	ret, err := Verify(pubKey, message, sign)
	if err != nil || !ret {
		t.Fatal("Verify with public key failed", err)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	origin := "Hello World"
	key := []byte(EncryptKey)
	encryptMsg, _ := Encrypt(key, []byte(origin))
	msg, err := Decrypt(key, []byte(encryptMsg))
	if err != nil {
		t.Fatal("Decrypt failed", err)
	}
	if string(msg) != origin {
		t.Fatal("Decrypt failed")
	}
}

func TestSerializeDeserializeKey(t *testing.T) {
	privKey, err := ToPrivKey(KeyString)
	if err != nil {
		t.Fatal("ToPrivKey failed", err)
	}
	privKeyString, err := FromPrivKey(privKey)
	if err != nil {
		t.Fatal("FromPrivKey failed", err)
	}
	if privKeyString != KeyString {
		t.Fatal("Serialize and deserialize private key failed", err)
	}

	pubKey := privKey.GetPublic()
	pubKeyString, err := FromPubKey(pubKey)
	if err != nil {
		t.Fatal("FromPubKey failed", err)
	}

	nPubKey, err := ToPubKey(pubKeyString)
	if err != nil {
		t.Fatal("ToPubKey failed", err)
	}

	pubkeyRaw, err := pubKey.Raw()
	if err != nil {
		t.Fatal("Get pubkey raw failed", err)
	}
	nPubkeyRaw, err := nPubKey.Raw()
	if err != nil {
		t.Fatal("Get PubKey raw failed", err)
	}

	if bytes.Compare(pubkeyRaw, nPubkeyRaw) != 0 {
		t.Fatal("Serialize and deserialize pub key failed", err)
	}
}

func TestHex64ToBase64(t *testing.T) {
	keyHex64 := "da146374a75310b9666e834ee4ad0866d6f4035967bfc76217c5a495fff9f0d0"
	privKey, err := Hex64ToBase64(keyHex64)
	if err != nil {
		t.Fatal("Decode hex64 private key failed", err)
	}

	priv, err := ToPrivKey(privKey)
	if err != nil {
		t.Fatal("Get private key failed", err)
	}
	pubKey := priv.GetPublic()
	pub, err := FromPubKey(pubKey)
	if err != nil {
		t.Fatal("From public key failed", err)
	}
	keyBase64 := "CAISIQJ/5o1cuJslw3ySQMIsbnMrvM/H/j5d3+N4rkNz48WCYw=="
	if pub != keyBase64 {
		t.Fatal("Base64 public key decode failed")
	}
}
