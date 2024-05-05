package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uint
}
type Service struct {
	secretKey string
	tokenExp  int
}

func NewAuthService(secretKey string, tokenExp int) *Service {
	return &Service{
		secretKey: secretKey,
		tokenExp:  tokenExp,
	}
}
func (s *Service) GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.tokenExp) * time.Hour * 24)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) ValidateToken(token string) (Claims, error) {
	claims := Claims{}
	if _, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	}); err != nil {
		return Claims{}, err
	}
	return claims, nil
}
