package lang

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

func computeHash(hash hash.Hash, filePath string) ([]byte, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	//hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}

	return hash.Sum(result), nil
}

func getHashingAlgo(algo string) hash.Hash {
	switch strings.ToLower(algo) {
	case "md5":
		return md5.New()
	case "sha256":
		return sha256.New()
	case "sha512":
		return sha512.New()
	case "sha1", "sha":
		return sha1.New()
	}
	return nil
}

func ValidateFile(algo string, hash string, filepath string) (err error) {

	hashing := getHashingAlgo(algo)

	if hashing == nil {
		return fmt.Errorf("Cannot compute hash with algorithm: %s", algo)
	}

	var cHash []byte
	if cHash, err = computeHash(hashing, filepath); err != nil {
		return err
	}

	hashBs, err := hex.DecodeString(hash)
	if err != nil {
		return err
	}
	if !bytes.Equal(hashBs, cHash) {
		return fmt.Errorf("Not equal: %s != %s", hash, cHash)
	}

	return nil
}
