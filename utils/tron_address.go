package utils

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	cr "github.com/tron-us/go-btfs-common/crypto"
	"github.com/tron-us/go-common/v2/crypto"
	"github.com/tron-us/protobuf/proto"

	eth "github.com/ethereum/go-ethereum/crypto"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

func GetTronPubKeyFromPubkey(pubkeyS string) (*string, error) {
	//pid, err := peer.IDB58Decode(peerId)
	//if err != nil {
	//	return nil, err
	//}
	pubkey, err := cr.ToPubKey(pubkeyS)
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

	addStr := hex.EncodeToString(addr.Bytes())
	result, err := crypto.Encode58Check(&addStr)

	return result, err
}

func GetTronPubKeyFromPeerIdPretty(peerId string) (*string, error) {
	pid, err := peer.IDB58Decode(peerId)
	if err != nil {
		return nil, err
	}
	pubkey, err := pid.ExtractPublicKey()
	if err != nil {
		return nil, err
	}

	return GetTronPubKeyFromIcPubKey(pubkey)
}

func TronSign(privKey ic.PrivKey, msg proto.Message) ([]byte, error) {
	raw, err := privKey.Raw()
	if err != nil {
		return nil, err
	}
	txBytes, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	ecdsa, err := eth.HexToECDSA(hex.EncodeToString(raw))
	sum := sha256.Sum256(txBytes)
	sig, err := eth.Sign(sum[:], ecdsa)
	return sig, err
}

func TronSignRaw(privKey ic.PrivKey, data []byte) ([]byte, error) {
	raw, err := privKey.Raw()
	if err != nil {
		return nil, err
	}
	ecdsa, err := eth.HexToECDSA(hex.EncodeToString(raw))
	sum := sha256.Sum256(data)
	sig, err := eth.Sign(sum[:], ecdsa)
	return sig, err
}

func GetTronPubKeyFromIcPubKey(pubkey ic.PubKey) (*string, error) {
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

	addStr := hex.EncodeToString(addr.Bytes())
	result, err := crypto.Encode58Check(&addStr)

	return result, err
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

func GetRawFullFromPeerIdPretty(peerid string) ([]byte, error) {
	peerId, err := peer.IDB58Decode(peerid)
	if err != nil {
		return nil, err
	}
	pubkey, err := peerId.ExtractPublicKey()
	if err != nil {
		return nil, err
	}
	return pubkey.Raw()
}

func GetPubKeyFromPeerId(pId string) (ic.PubKey, error) {
	peerId, err := peer.IDB58Decode(pId)
	if err != nil {
		return nil, err
	}
	pubKey, err2 := peerId.ExtractPublicKey()
	if err2 != nil {
		return nil, err2
	}
	return pubKey, nil
}
