package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// HashStudentID 对学号进行 SHA-256 哈希（确定性，用于数据库查找）
func HashStudentID(studentID string) string {
	h := sha256.Sum256([]byte(studentID))
	return hex.EncodeToString(h[:])
}

// Encrypt AES-256-CBC 加密（用于学号等敏感字段）
func Encrypt(plaintext string, key string) (string, error) {
	// 取 key 前 32 字节作为 AES-256 密钥
	keyBytes := []byte(key)
	if len(keyBytes) > 32 {
		keyBytes = keyBytes[:32]
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("创建加密器失败: %w", err)
	}

	plainBytes := []byte(plaintext)
	// PKCS7 填充
	blockSize := block.BlockSize()
	padding := blockSize - len(plainBytes)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	plainBytes = append(plainBytes, padText...)

	// CBC 模式加密
	ciphertext := make([]byte, aes.BlockSize+len(plainBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("生成 IV 失败: %w", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plainBytes)

	return hex.EncodeToString(ciphertext), nil
}

// Decrypt AES-256-CBC 解密
func Decrypt(ciphertext string, key string) (string, error) {
	keyBytes := []byte(key)
	if len(keyBytes) > 32 {
		keyBytes = keyBytes[:32]
	}

	data, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("解码密文失败: %w", err)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("创建解密器失败: %w", err)
	}

	if len(data) < aes.BlockSize {
		return "", fmt.Errorf("密文数据过短")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)

	// PKCS7 去填充
	padding := int(data[len(data)-1])
	if padding > len(data) || padding > aes.BlockSize {
		return "", fmt.Errorf("填充数据无效")
	}
	for i := 0; i < padding; i++ {
		if data[len(data)-1-i] != byte(padding) {
			return "", fmt.Errorf("填充校验失败")
		}
	}

	return string(data[:len(data)-padding]), nil
}
