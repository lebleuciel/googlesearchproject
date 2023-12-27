package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveEncryptedFile(fileContent []byte, destDir string, key []byte) (string, error) {
	// Generate UUID for file name
	newFileName, err := GenerateUUID()
	if err != nil {
		return "", err
	}

	cipherText, err := EncryptFile(fileContent, key)
	if err != nil {
		return "", err
	}

	// Specify the directory where you want to save the file
	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Specify the path for the encrypted file
	destPath := filepath.Join(destDir, newFileName)

	err = os.WriteFile(destPath, cipherText, os.ModePerm)
	if err != nil {
		return "", err
	}

	return newFileName, nil
}

func ReadFileContent(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	content, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func EncryptFile(plainText []byte, key []byte) ([]byte, error) {
	// Creating block of algorithm
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Creating GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generating random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Decrypt file
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)
	return cipherText, nil
}

func DecryptFile(filePath string, key []byte) error {
	cipherText, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Creating block of algorithm
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Creating GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Deattached nonce and decrypt
	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, plainText, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
