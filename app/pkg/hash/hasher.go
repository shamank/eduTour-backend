package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// PasswordHasher provides hashing logic to securely store passwords.
type PasswordHasher interface {
	Hash(str string) (string, error)
	SimpleHash(str string) (string, error)
}

// SHA256Hasher uses SHA1 to hash passwords with provided salt.
type SHA256Hasher struct {
	salt string
}

func NewSHA256Hasher(salt string) *SHA256Hasher {
	return &SHA256Hasher{salt: salt}
}

// Hash creates SHA256 hash of given password.
func (h *SHA256Hasher) Hash(str string) (string, error) {
	hash := sha256.New()

	if _, err := hash.Write([]byte(str)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}

func (h *SHA256Hasher) SimpleHash(str string) (string, error) {
	// Создание нового хэш-объекта SHA-256
	hash := sha256.New()

	// Преобразование пароля в байтовый массив и запись его в хэш-объект
	_, err := hash.Write([]byte(str))
	if err != nil {
		return "", err
	}

	// Получение хэш-суммы в виде байтового массива
	hashBytes := hash.Sum(nil)

	// Преобразование хэш-суммы в строку в шестнадцатеричном формате
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}
