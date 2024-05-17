package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/mrkovshik/yandex_diploma/internal/auth"
)

func (s *restAPIServer) Auth(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
			return
		}
		claims := &auth.Claims{}
		if _, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(s.cfg.SecretKey), nil
		}); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
			return
		}
		if _, err := s.storage.GetUserByID(ctx, claims.UserID); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
			return
		}

		c.Set("userID", claims.UserID)
	}
}

//TODO: add compressing MW here
