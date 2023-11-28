package helpers

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func Encrypt(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(uuid), nil
}

func SplitBySpaceComma(input []string) []string {
	var result []string

	// Iterate through each tag and split by space and comma
	for _, s := range input {
		// Split by space and comma
		parts := strings.FieldsFunc(s, func(r rune) bool {
			return r == ' ' || r == ','
		})

		// Add the split parts to the overall slice
		result = append(result, parts...)
	}
	return result
}
