package rest

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mrkovshik/yandex_diploma/internal/auth"
	"github.com/mrkovshik/yandex_diploma/internal/storage/postgres"
)

func (s restApiServer) Auth(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, found := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
		if !found {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		claims := &auth.Claims{}
		if _, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(s.cfg.SecretKey), nil
		}); err != nil {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		if _, err := postgres.NewPostgresUserStorage(s.db).GetUserByID(ctx, claims.UserID); err != nil {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
	}
}
