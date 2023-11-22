package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gammazero/workerpool"
	"github.com/joho/godotenv"
)

type workerPoolClient struct {
	instance *workerpool.WorkerPool
}

func NewWorkerPoolClient() *workerPoolClient {
	return &workerPoolClient{nil}
}

func (instance *workerPoolClient) Instance() (*workerpool.WorkerPool, error) {
	if instance.instance == nil {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Error while load .env file: " + err.Error())
			return nil, ErrLoadGormEnvFile
		} else {
			wpSize, _ := strconv.Atoi(os.Getenv("WORKER_POOL_SIZE"))
			instance.instance = workerpool.New(wpSize)
		}
	}
	return instance.instance, nil
}
