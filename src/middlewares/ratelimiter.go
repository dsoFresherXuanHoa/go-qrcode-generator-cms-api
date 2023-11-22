package middlewares

import (
	"errors"
	"go-qrcode-generator-cms-api/src/configs"
	"go-qrcode-generator-cms-api/src/constants"
	"go-qrcode-generator-cms-api/src/entity"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	ErrTooManyRequest = errors.New("you make too many request, try again after 5 minutes")
)

func RateLimit() gin.HandlerFunc {
	refreshIpsListTime, _ := strconv.Atoi(os.Getenv("REFRESH_IPS_LIST_AFTER"))

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Second * 30)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > time.Duration(refreshIpsListTime)*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(ctx *gin.Context) {
		limit, _ := configs.NewRateLimitClient().Instance()
		ip := ctx.ClientIP()
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: limit}
		}
		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, entity.NewStandardResponse(nil, http.StatusTooManyRequests, constants.StatusTooManyRequests, ErrTooManyRequest.Error(), ErrTooManyRequest.Error()))
		} else {
			mu.Unlock()
			ctx.Next()
		}
	}
}
