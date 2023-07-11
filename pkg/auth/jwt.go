package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type userClaims struct {
	UserID   int
	UserName string
	Role     string
	ExpireAt int64
}

type TokenManager interface {
	Generate(userID int, userName string, role string) (string, time.Duration, error)
	Parse(token string) (userClaims, error)
	GenerateToken(byteSize int) (string, error)
	GenerateRefreshToken() (string, int64, error)
}

type Manager struct {
	signedKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewManager(signedKey string, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) (*Manager, error) {
	if signedKey == "" {
		return nil, errors.New("empty signed key")
	}

	return &Manager{
		signedKey:       signedKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}, nil
}

func (m *Manager) Generate(userID int, userName string, role string) (string, time.Duration, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userID,
		"user_role": role,
		"user_name": userName,
		"expire_at": time.Now().Add(m.accessTokenTTL).Unix(),
	})
	tokenString, err := token.SignedString([]byte(m.signedKey))

	if err != nil {
		return "", 0, fmt.Errorf("error with sign token: %s", err.Error())
	}

	return tokenString, m.accessTokenTTL, nil
}

func (m *Manager) Parse(tokenString string) (userClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signedKey), nil
	})
	if err != nil {
		return userClaims{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return userClaims{
			UserID:   int(claims["user_id"].(float64)),
			UserName: claims["user_name"].(string),
			Role:     claims["user_role"].(string),
			ExpireAt: int64(claims["expire_at"].(float64)),
		}, nil
	}
	return userClaims{}, fmt.Errorf("cannot get claims from token")

}

func (m *Manager) GenerateToken(byteSize int) (string, error) {
	bytes := make([]byte, byteSize)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Кодируем байты в строку
	token := base64.RawURLEncoding.EncodeToString(bytes)
	return token, nil
}

func (m *Manager) GenerateRefreshToken() (string, int64, error) {
	token, err := m.GenerateToken(32)
	if err != nil {
		return "", 0, err
	}
	expireAt := time.Now().Add(m.refreshTokenTTL).Unix()
	return token, expireAt, nil
}
