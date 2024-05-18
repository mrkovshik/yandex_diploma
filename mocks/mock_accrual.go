package mock_api

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mrkovshik/yandex_diploma/internal/config"
	"github.com/mrkovshik/yandex_diploma/internal/model"
)

const (
	NumberForNotFound        = "1234567890"
	NumberForInternalErr     = "9876543210"
	NumberForTooManyRequests = "2468013579"
)

var states = []model.AccrualState{
	model.AccrualStateProcessing,
	model.AccrualStateRegistered,
	model.AccrualStateInvalid,
	model.AccrualStateProcessed,
}

func Run(cfg *config.Config) {
	r := gin.Default()
	r.GET("api/orders/:order", func(c *gin.Context) {
		var resp model.AccrualResponse
		if err := c.ShouldBindUri(&resp); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		switch resp.Order {
		case NumberForNotFound:
			c.AbortWithStatus(http.StatusNoContent)
			return
		case NumberForInternalErr:
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		case NumberForTooManyRequests:
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		default:
			rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
			resp.Status = states[rand.Intn(len(states))]
			if resp.Status == model.AccrualStateProcessed {
				resp.Accrual = rand.Intn(10000)
			}
			c.JSON(http.StatusOK, resp)
		}

	})
	err := r.Run(cfg.AccrualSystemAddress)
	log.Fatal(err)

}
