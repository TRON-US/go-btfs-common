package crypto

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcutil/base58"
	ethCommon "github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/crypto"
	ic "github.com/libp2p/go-libp2p/core/crypto"
)

const (
	AddressLength = 21
	AddressPrefix = "41"
)

type Address [AddressLength]byte

type Keys struct {
	Base58Address string
	HexAddress    string
	HexPrivateKey string
	HexPubKey     string
	Base64PubKey  string
}

func FromIcPrivateKey(privKey ic.PrivKey) (*Keys, error) {
	keys := &Keys{}
	pubKey := privKey.GetPublic()
	var err error
	keys.Base64PubKey, err = FromPubKey(pubKey)
	if err != nil {
		return nil, err
	}

	pubKeyRaw, err := Secp256k1PublicKeyRaw(pubKey)
	if err != nil {
		return nil, err
	}
	keys.HexPubKey = hex.EncodeToString(pubKeyRaw)

	privKeyRaw, err := privKey.Raw()
	if err != nil {
		return nil, err
	}
	keys.HexPrivateKey = hex.EncodeToString(privKeyRaw)

	// test for exchange address
	privateKey, err := eth.HexToECDSA(keys.HexPrivateKey)
	if err != nil {
		return nil, err
	}
	if privateKey == nil {
		return nil, err
	}
	addr, err := PublicKeyToAddress(privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	keys.HexAddress = hex.EncodeToString(addr.Bytes())
	keys.Base58Address, err = Encode58Check(addr.Bytes())
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func FromPrivateKey(key string) (*Keys, error) {
	privKey, err := ToPrivKey(key)
	if err != nil {
		priv_key, err := Hex64ToBase64(key)
		if err != nil {
			return nil, err
		}
		privKey, err = ToPrivKey(priv_key)
		if err != nil {
			return nil, err
		}
	}
	return FromIcPrivateKey(privKey)
}

// ecdsa key to Tron address
func PublicKeyToAddress(p ecdsa.PublicKey) (Address, error) {
	addr := eth.PubkeyToAddress(p)

	addressTron := make([]byte, AddressLength)

	addressPrefix, err := FromHex(AddressPrefix)
	if err != nil {
		return Address{}, err
	}

	addressTron = append(addressTron, addressPrefix...)
	addressTron = append(addressTron, addr.Bytes()...)

	return BytesToAddress(addressTron), nil
}

func (a *Address) Bytes() []byte {
	return a[:]
}

func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// Convert byte to address.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// Decode hex string as bytes
func FromHex(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, errors.New("empty hex string")
	}

	return hex.DecodeString(input[:])
}

// Decode by base58 and check.
func Decode58Check(input string) ([]byte, error) {
	decodeCheck := base58.Decode(input)
	if len(decodeCheck) <= 4 {
		return nil, errors.New("base58 encode length error")
	}
	decodeData := decodeCheck[:len(decodeCheck)-4]
	hash0, err := Hash(decodeData)
	if err != nil {
		return nil, err
	}
	hash1, err := Hash(hash0)
	if hash1 == nil {
		return nil, err
	}
	if hash1[0] == decodeCheck[len(decodeData)] && hash1[1] == decodeCheck[len(decodeData)+1] &&
		hash1[2] == decodeCheck[len(decodeData)+2] && hash1[3] == decodeCheck[len(decodeData)+3] {
		return decodeData, nil
	}
	return nil, errors.New("base58 check failed")
}

// Encode by base58 and check.
func Encode58Check(input []byte) (string, error) {
	h0, err := Hash(input)
	if err != nil {
		return "", err
	}
	h1, err := Hash(h0)
	if err != nil {
		return "", err
	}
	if len(h1) < 4 {
		return "", errors.New("base58 encode length error")
	}
	inputCheck := append(input, h1[:4]...)

	return base58.Encode(inputCheck), nil
}

// Package goLang sha256 hash algorithm.
func Hash(s []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(s)
	if err != nil {
		return nil, err
	}
	bs := h.Sum(nil)
	return bs, nil
}

// Get Tron address from ledger address
func AddressLedgerToTron(ledgerAddress []byte) (Address, error) {
	addr := ethCommon.BytesToAddress(eth.Keccak256(ledgerAddress[1:])[12:])
	addressTron := make([]byte, AddressLength)
	addressPrefix, err := FromHex(AddressPrefix)
	if err != nil {
		return Address{}, err
	}
	addressTron = append(addressTron, addressPrefix...)
	addressTron = append(addressTron, addr.Bytes()...)
	return BytesToAddress(addressTron), nil
}
