package encoding

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"strings"
)

// Base64Encode returns the base64 encoding of the input data
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode decodes a base64 string into bytes
func Base64Decode(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// Base64URLEncode returns the URL-safe base64 encoding of the input data
func Base64URLEncode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

// Base64URLDecode decodes a URL-safe base64 string into bytes
func Base64URLDecode(encoded string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(encoded)
}

// Base32Encode returns the base32 encoding of the input data
func Base32Encode(data []byte) string {
	return base32.StdEncoding.EncodeToString(data)
}

// Base32Decode decodes a base32 string into bytes
func Base32Decode(encoded string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(encoded)
}

// HexEncode returns the hexadecimal encoding of the input data
func HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

// HexDecode decodes a hexadecimal string into bytes
func HexDecode(encoded string) ([]byte, error) {
	return hex.DecodeString(encoded)
}

// StringToBase64 encodes a string to base64
func StringToBase64(s string) string {
	return Base64Encode([]byte(s))
}

// Base64ToString decodes a base64 string back to a regular string
func Base64ToString(encoded string) (string, error) {
	decoded, err := Base64Decode(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// JSONMarshal marshals an object to JSON with indentation
func JSONMarshal(v interface{}, indent bool) ([]byte, error) {
	if indent {
		return json.MarshalIndent(v, "", "  ")
	}
	return json.Marshal(v)
}

// JSONUnmarshal unmarshals JSON data into an object
func JSONUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// XMLMarshal marshals an object to XML with indentation
func XMLMarshal(v interface{}, indent bool) ([]byte, error) {
	if indent {
		return xml.MarshalIndent(v, "", "  ")
	}
	return xml.Marshal(v)
}

// XMLUnmarshal unmarshals XML data into an object
func XMLUnmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

// EncodeToFile encodes data and writes it to a file
func EncodeToFile(filename string, data []byte, encodingFunc func([]byte) string) error {
	encoded := encodingFunc(data)
	return ioutil.WriteFile(filename, []byte(encoded), 0644)
}

// DecodeFromFile reads encoded data from a file and decodes it
func DecodeFromFile(filename string, decodingFunc func(string) ([]byte, error)) ([]byte, error) {
	encodedData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return decodingFunc(string(encodedData))
}

// ConvertEncoding converts data from one encoding to another
func ConvertEncoding(data string, decodeFunc func(string) ([]byte, error), encodeFunc func([]byte) string) (string, error) {
	decoded, err := decodeFunc(data)
	if err != nil {
		return "", err
	}
	return encodeFunc(decoded), nil
}

// IsBase64 checks if a string is valid base64
func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

// IsHex checks if a string is valid hexadecimal
func IsHex(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}

// RemoveWhitespace removes all whitespace from a string
func RemoveWhitespace(s string) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
			return -1
		}
		return r
	}, s)
}
