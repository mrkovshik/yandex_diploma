package rest

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_getOrderNumberFromContext(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name    string
		number  string
		want    uint
		errWant bool
	}{
		{"luhn_invalid", "1234", 0, true},
		{"int_invalid", "d234", 0, true},
		{"empty", "", 0, true},
		{"valid", "12345678903", 12345678903, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			reqBody := bytes.NewBufferString(tt.number)
			c.Request, _ = http.NewRequest(http.MethodPost, "/mock", reqBody)
			orderNumber, err := getOrderNumberFromContext(c)
			assert.Equal(t, !errors.Is(err, nil), tt.errWant)
			assert.Equal(t, tt.want, orderNumber)
		})
	}
}
