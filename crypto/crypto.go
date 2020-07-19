package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/gogo/protobuf/proto"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	pb "github.com/libp2p/go-libp2p-core/crypto/pb"
	"github.com/libp2p/go-libp2p-core/peer"
)

var (
	secp256k1N, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
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
func ToPubKey(pubKey string) (ic.PubKey, error) {
	raw, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	return ic.UnmarshalPublicKey(raw)
}

func FromPubKey(pubKey ic.PubKey) (string, error) {
	pkb, err := ic.MarshalPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(pkb), nil
}

func FromPrivKey(privKey ic.PrivKey) (string, error) {
	prkb, err := ic.MarshalPrivateKey(privKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(prkb), nil
}

// Secp256k1 private key string to ic.PrivKey interface
func ToPrivKeyRaw(privKey []byte) (ic.PrivKey, error) {
	return ic.UnmarshalSecp256k1PrivateKey(privKey)
}

// public key string to ic.PubKey interface
func ToPubKeyRaw(pubKey []byte) (ic.PubKey, error) {
	return ic.UnmarshalSecp256k1PublicKey(pubKey)
}

func GenKeyPairs() (ic.PrivKey, ic.PubKey, error) {
	return ic.GenerateSecp256k1Key(rand.Reader)
}

func addBase64Padding(text []byte) string {
	value := string(text)
	m := len(value) % 4
	if m != 0 {
		value += strings.Repeat("=", 4-m)
	}

	return value
}

func removeBase64Padding(value string) string {
	return strings.Replace(value, "=", "", -1)
}

// pkcs7 padding
func Pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func Unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}

func Encrypt(key, text []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	msg := Pad(text)
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(msg))
	finalMsg := removeBase64Padding(base64.URLEncoding.EncodeToString(ciphertext))
	return finalMsg, nil
}

func Decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	decodedMsg, err := base64.URLEncoding.DecodeString(addBase64Padding(text))
	if err != nil {
		return nil, err
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return nil, errors.New("blocksize must be multiple of decoded message length")
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	unpadMsg, err := Unpad(msg)
	if err != nil {
		return nil, err
	}

	return unpadMsg, nil
}

func GetPubKeyFromPeerId(pid string) (ic.PubKey, error) {
	peerId, err := peer.IDB58Decode(pid)
	if err != nil {
		return nil, err
	}
	pubKey, err2 := peerId.ExtractPublicKey()
	if err2 != nil {
		return nil, err2
	}
	return pubKey, nil
}

func Hex64ToBase64(key string) (string, error) {
	src := []byte(key)

	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	if err != nil {
		return "", fmt.Errorf("decode hex64 failed: %v", err)
	}

	// marshal
	pbmes := new(pb.PrivateKey)
	pbmes.Type = pb.KeyType_Secp256k1
	pbmes.Data = dst
	marshaledKey, err := proto.Marshal(pbmes)
	if err != nil {
		return "", fmt.Errorf("marshal key failed: %v", err)
	}

	// base64 encoding
	encodeKey := base64.StdEncoding.EncodeToString(marshaledKey)
	return encodeKey, nil
}

// GetPrivKeyFromHexOrBase64 can decode a priv key from either hex or base64
// format to satisfy different key storage encoding schemes
func GetPrivKeyFromHexOrBase64(raw string) (ic.PrivKey, error) {
	key, err := Hex64ToBase64(raw)
	if err != nil {
		// Check base64 format directly
		key = raw
	}
	return ToPrivKey(key)
}

// HexToECDSA parses a secp256k1 private key.
func HexToECDSA(hexkey string) (*ecdsa.PrivateKey, error) {
	b, err := hex.DecodeString(hexkey)
	if byteErr, ok := err.(hex.InvalidByteError); ok {
		return nil, fmt.Errorf("invalid hex character %q in private key", byte(byteErr))
	} else if err != nil {
		return nil, errors.New("invalid hex data for private key")
	}
	return ToECDSA(b)
}

// ToECDSA creates a private key with the given D value.
func ToECDSA(d []byte) (*ecdsa.PrivateKey, error) {
	return toECDSA(d, true)
}

// toECDSA creates a private key with the given D value. The strict parameter
// controls whether the key's length should be enforced at the curve size or
// it can also accept legacy encodings (0 prefixes).
func toECDSA(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P256()
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N
	if priv.D.Cmp(secp256k1N) >= 0 {
		return nil, fmt.Errorf("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}
	return priv, nil
}
