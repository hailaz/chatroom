package service

import (
	"context"
	"fmt"
	"goframechat/internal/consts"
	"goframechat/internal/model/entity"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v4"
)

// JwtService handles JWT token generation and validation
type JwtService struct {
	secretKey []byte
}

// JwtClaims represents the custom JWT claims
type JwtClaims struct {
	UserId   uint   `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewJwtService creates a new JwtService instance
func NewJwtService() *JwtService {
	secretKey := []byte(g.Cfg().MustGet(context.Background(), "jwt.secretKey", "goframechat_secret_key").String())
	return &JwtService{
		secretKey: secretKey,
	}
}

// GenerateToken generates a new JWT token for a user
func (s *JwtService) GenerateToken(user *entity.User) (string, error) {
	// Create claims
	claims := JwtClaims{
		UserId:   user.Id,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(consts.JwtExpireTime) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    consts.JwtIssuer,
		},
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ParseToken parses a JWT token string and returns the claims
func (s *JwtService) ParseToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
