package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenManager interface {
	Generate(userID string, roleID string, ttl time.Duration) (string, error)
	Parse(token string) (string, string, error)
	GenerateRefreshToken() (string, error)
}

type Manager struct {
	signedKey string
}

func NewManager(signedKey string) (*Manager, error) {
	if signedKey == "" {
		return nil, errors.New("empty signed key")
	}

	return &Manager{signedKey: signedKey}, nil
}

func (m *Manager) Generate(userID string, roleID string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userID,
		"user_role": roleID,
		"expire_at": time.Now().Add(ttl).Unix(),
	})
	tokenString, err := token.SignedString([]byte(m.signedKey))
	if err != nil {
		return "", fmt.Errorf("error with sign token: %s", err.Error())
	}
	return tokenString, nil
}

func (m *Manager) Parse(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return m.signedKey, nil
	})
	if err != nil {
		return "", "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user_id"].(string), claims["user_role"].(string), nil
	} else {
		return "", "", fmt.Errorf("cannot get claims from token")
	}
}

func (m *Manager) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Кодируем байты в строку
	token := base64.StdEncoding.EncodeToString(bytes)
	return token, nil
}
