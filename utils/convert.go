package utils

import (
	"encoding/base64"
	"fmt"
)

const (
	Text = iota + 1
	Base64
)

func StringToBytes(str string, encoding int) ([]byte, error) {
	switch encoding {
	case Text:
		return []byte(str), nil
	case Base64:
		by, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			return nil, err
		}
		return by, nil
	default:
		return nil, fmt.Errorf(`unexpected encoding [%d], expected 1(Text) or 2(Base64)`, encoding)
	}
}

func BytesToString(data []byte, encoding int) (string, error) {
	switch encoding {
	case Text:
		return string(data), nil
	case Base64:
		return base64.StdEncoding.EncodeToString(data), nil
	default:
		return "", fmt.Errorf(`unexpected parameter [%d] is given, either "text" or "base64" should be used`, encoding)
	}
}
