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
	Roles    []string
	ExpireAt int64
}

type TokenManager interface {
	Generate(userID int, userName string, roles []string) (string, time.Duration, error)
	Parse(token string) (userClaims, error)
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

func (m *Manager) Generate(userID int, userName string, roles []string) (string, time.Duration, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"user_roles": roles,
		"user_name":  userName,
		"expire_at":  time.Now().Add(m.accessTokenTTL).Unix(),
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
		roles := make([]string, len(claims["user_roles"].([]interface{})))
		for i, v := range claims["user_roles"].([]interface{}) {
			roles[i], ok = v.(string)
			if !ok {
				return userClaims{}, errors.New("error occurred parsing claims (user roles)")
			}
		}
		return userClaims{
			UserID:   int(claims["user_id"].(float64)),
			UserName: claims["user_name"].(string),
			Roles:    roles,
			ExpireAt: int64(claims["expire_at"].(float64)),
		}, nil
	}
	return userClaims{}, fmt.Errorf("cannot get claims from token")

}

func (m *Manager) GenerateRefreshToken() (string, int64, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", 0, err
	}

	// Кодируем байты в строку
	token := base64.StdEncoding.EncodeToString(bytes)
	expireAt := time.Now().Add(m.refreshTokenTTL).Unix()
	return token, expireAt, nil
}
