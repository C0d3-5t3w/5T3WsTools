package hashExt

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"strings"
)

// StringToMD5 returns MD5 hash of the input string
func StringToMD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// StringToSHA1 returns SHA1 hash of the input string
func StringToSHA1(text string) string {
	hasher := sha1.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// StringToSHA256 returns SHA256 hash of the input string
func StringToSHA256(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// StringToSHA512 returns SHA512 hash of the input string
func StringToSHA512(text string) string {
	hasher := sha512.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// StringToCRC32 returns CRC32 checksum of the input string
func StringToCRC32(text string) string {
	hasher := crc32.NewIEEE()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// FileToHash returns the hash of a file using the provided hash algorithm
func FileToHash(filepath string, hasher hash.Hash) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// FileToMD5 returns MD5 hash of a file
func FileToMD5(filepath string) (string, error) {
	return FileToHash(filepath, md5.New())
}

// FileToSHA1 returns SHA1 hash of a file
func FileToSHA1(filepath string) (string, error) {
	return FileToHash(filepath, sha1.New())
}

// FileToSHA256 returns SHA256 hash of a file
func FileToSHA256(filepath string) (string, error) {
	return FileToHash(filepath, sha256.New())
}

// FileToSHA512 returns SHA512 hash of a file
func FileToSHA512(filepath string) (string, error) {
	return FileToHash(filepath, sha512.New())
}

// FileToCRC32 returns CRC32 checksum of a file
func FileToCRC32(filepath string) (string, error) {
	return FileToHash(filepath, crc32.NewIEEE())
}

// CompareHashes compares two hashes (case-insensitive)
func CompareHashes(hash1, hash2 string) bool {
	return strings.EqualFold(hash1, hash2)
}
