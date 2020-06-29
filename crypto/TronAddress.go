package crypto

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcutil/base58"
	eth "github.com/ethereum/go-ethereum/crypto"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

func GetTronPubKeyFromPeerIdPretty(peerId string) (*string, error) {
	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return nil, err
	}
	pubkey, err := pid.ExtractPublicKey()
	if err != nil {
		return nil, err
	}

	pubkeyRaw, err := ic.RawFull(pubkey)
	if err != nil {
		return nil, err
	}

	ethPubkey, err := eth.UnmarshalPubkey(pubkeyRaw)
	if err != nil {
		return nil, err
	}

	addr, err := EcdsaPublicKeyToAddress(*ethPubkey)
	if err != nil {
		return nil, err
	}
	result, err := Encode58CheckBytes(addr.Bytes())
	return &result, err
}
func EcdsaPublicKeyToAddress(p ecdsa.PublicKey) (TronAddress, error) {
	addr := eth.PubkeyToAddress(p)

	addressTron := make([]byte, AddressLength)

	addressPrefix, err := FromHex(AddressPrefix)
	if err != nil {
		return TronAddress{}, err
	}

	addressTron = append(addressTron, addressPrefix...)
	addressTron = append(addressTron, addr.Bytes()...)

	return BytesToAddress(addressTron), nil
}

func sha256_hash(s []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(s)
	if err != nil {
		return nil, err
	}
	bs := h.Sum(nil)
	return bs, nil
}

func Encode58CheckBytes(input []byte) (string, error) {
	h0, err := sha256_hash(input)
	if err != nil {
		return "", err
	}
	h1, err := sha256_hash(h0)
	if err != nil {
		return "", err
	}
	if len(h1) < 4 {
		return "", errors.New("base58 encode length error")
	}
	inputCheck := append(input, h1[:4]...)

	return base58.Encode(inputCheck), nil
}

const (
	AddressLength = 21
	AddressPrefix = "41"
)

type TronAddress [AddressLength]byte

func BytesToAddress(b []byte) TronAddress {
	var a TronAddress
	a.SetBytes(b)
	return a
}

func (a *TronAddress) Bytes() []byte {
	return a[:]
}

// Convert byte to address.
func (a *TronAddress) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

func FromHex(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, errors.New("empty hex string")
	}

	return hex.DecodeString(input[:])
}
