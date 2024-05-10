package accrual

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrkovshik/yandex_diploma/internal/app_errors"
	"github.com/stretchr/testify/assert"
)

func Test_service_GetOrderScore(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/orders/123":
			mockResp := Response{
				Order:   123,
				Status:  "PROCESSING",
				Accrual: 0,
			}
			respJSON, _ := json.Marshal(mockResp)
			w.WriteHeader(http.StatusOK)
			w.Write(respJSON)

		case "/api/orders/456":
			w.WriteHeader(http.StatusNoContent)
		case "/api/orders/789":
			w.WriteHeader(http.StatusTooManyRequests)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer mockServer.Close()
	s := NewAccrualService(mockServer.URL[7:])
	tests := []struct {
		name        string
		orderNumber uint
		want        Response
		err         error
	}{
		{"1_pos", 123, Response{
			Order:   123,
			Status:  "PROCESSING",
			Accrual: 0,
		}, nil},
		{"2_neg", 456, Response{
			Order:   0,
			Status:  "",
			Accrual: 0,
		}, app_errors.ErrNoSuchOrder},
		{"3_neg", 789, Response{
			Order:   0,
			Status:  "",
			Accrual: 0,
		}, app_errors.ErrTooManyRetrials},
		{"4_neg", 7819, Response{
			Order:   0,
			Status:  "",
			Accrual: 0,
		}, app_errors.ErrInvalidResponseCode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetOrderAccrual(tt.orderNumber)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, errors.Is(err, tt.err), true)
		})
	}
}
