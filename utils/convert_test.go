package utils

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/tron-us/go-btfs-common/protos/escrow"
	"github.com/tron-us/protobuf/proto"
)

func TestBytesToStringAndStringToBytesInText(t *testing.T) {
	data := []byte("Hello, welcome to this BytesToString/StringToBytes unit test")

	str, err := BytesToString(data, Text)
	if err != nil {
		t.Fatal(err)
	}

	databack, err := StringToBytes(str, Text)
	if err != nil {
		t.Fatal(err)
	}

	if res := bytes.Compare(data, databack); res != 0 {
		t.Fatal("original bytes and converted back bytes don't match")
	}
}

func TestBytesToStringAndStringToBytesInBase64(t *testing.T) {
	escrowContract := new(escrow.EscrowContract)
	escrowContract.Amount = 101.00
	contract := &escrow.SignedEscrowContract{Contract: escrowContract}
	contractBytes, err := proto.Marshal(contract)
	if err != nil {
		t.Fatal(err)
	}

	priv, _, err := crypto.GenerateSecp256k1Key(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	sig, err := priv.Sign(contractBytes)
	if err != nil {
		t.Fatal(err)
	}

	str, err := BytesToString(sig, Base64)
	if err != nil {
		t.Fatal(err)
	}

	sigback, err := StringToBytes(str, Base64)
	if err != nil {
		t.Fatal(err)
	}

	if res := bytes.Compare(sig, sigback); res != 0 {
		t.Fatal("original bytes and converted back bytes don't match")
	}
}
