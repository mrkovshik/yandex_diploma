package accrual

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mrkovshik/yandex_diploma/internal/apperrors"
	"github.com/mrkovshik/yandex_diploma/internal/model"
)

func Test_service_GetOrderScore(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/orders/123":
			mockResp := model.AccrualResponse{
				Order:   "123",
				Status:  "PROCESSING",
				Accrual: 0,
			}
			respJSON, _ := json.Marshal(mockResp)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(respJSON)

		case "/api/orders/456":
			w.WriteHeader(http.StatusNoContent)
		case "/api/orders/789":
			w.WriteHeader(http.StatusTooManyRequests)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer mockServer.Close()
	s := NewAccrualService(mockServer.URL)
	tests := []struct {
		name        string
		orderNumber string
		want        model.AccrualResponse
		err         error
	}{
		{"1_pos", "123", model.AccrualResponse{
			Order:   "123",
			Status:  "PROCESSING",
			Accrual: 0,
		}, nil},
		{"2_neg", "456", model.AccrualResponse{
			Order:   "",
			Status:  "",
			Accrual: 0,
		}, apperrors.ErrNoSuchOrder},
		{"3_neg", "789", model.AccrualResponse{
			Order:   "",
			Status:  "",
			Accrual: 0,
		}, apperrors.ErrTooManyRetrials},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetOrderAccrual(tt.orderNumber)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, errors.Is(err, tt.err), true)
		})
	}
}
