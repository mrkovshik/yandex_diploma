package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	key = "akjsfdsf"
	exp = 3
)

func TestService_GenerateToken(t *testing.T) {

	tests := []struct {
		name   string
		userID uint
	}{
		{"1", 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				secretKey: key,
				tokenExp:  exp,
			}
			token, err := s.GenerateToken(tt.userID)
			assert.NoError(t, err)
			claims, err1 := s.ValidateToken(token)
			assert.NoError(t, err1)
			assert.Equal(t, tt.userID, claims.UserID)

		})
	}
}
