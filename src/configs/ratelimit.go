package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

var ()

type rateLimitClient struct {
	instance *rate.Limiter
}

func NewRateLimitClient() *rateLimitClient {
	return &rateLimitClient{nil}
}

func (instance *rateLimitClient) Instance() (*rate.Limiter, error) {
	if instance.instance == nil {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Error while load .env file: " + err.Error())
			return nil, ErrLoadGormEnvFile
		} else {
			var maxRequestAllow, _ = strconv.Atoi(os.Getenv("MAX_REQUEST_ALLOW_RATE_LIMIT"))
			resetRateLimitAfter, _ := strconv.Atoi(os.Getenv("RESET_RATE_LIMIT_AFTER"))
			r := rate.Every(time.Duration(resetRateLimitAfter) * time.Minute)
			instance.instance = rate.NewLimiter(r, maxRequestAllow)
		}
	}
	return instance.instance, nil
}
