package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"

	"github.com/gogo/protobuf/proto"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
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
	return ic.UnmarshalSecp256k1PublicKey(raw)
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
