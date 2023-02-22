package crypto

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"

	"github.com/tron-us/go-common/v2/crypto"
	"github.com/tron-us/protobuf/proto"

	eth "github.com/ethereum/go-ethereum/crypto"
	ic "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

func GetTronPubKeyFromPubkey(pubkeyS string) (*string, error) {
	pubkey, err := ToPubKey(pubkeyS)
	if err != nil {
		return nil, err
	}

	pubkeyRaw, err := ic.MarshalPublicKey(pubkey)
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
	pid, err := peer.IDFromBytes([]byte(peerId))
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
	pubkeyRaw, err := ic.MarshalPublicKey(pubkey)
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
func EcdsaPublicKeyToAddress(p ecdsa.PublicKey) (Address, error) {
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

func GetRawFullFromPeerIdPretty(peerid string) ([]byte, error) {
	peerId, err := peer.IDFromBytes([]byte(peerid))
	if err != nil {
		return nil, err
	}
	pubkey, err := peerId.ExtractPublicKey()
	if err != nil {
		return nil, err
	}
	return pubkey.Raw()
}
