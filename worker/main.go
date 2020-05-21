package worker

import (
	"sync"
	"time"

	"github.com/sjatkins12/workerServer/cache"
	"github.com/sjatkins12/workerServer/utils"
)

const interval time.Duration = 30

// Job ... Represents a request to process a task by the worker
type Job struct {
	Response chan JobResponse
	Task     TaskType
}

// JobResponse ... Represents the response the worker will reply with
type JobResponse struct {
	Data interface{}
	Err  error
}

// Worker ... handle dispatching tasks to be processed
func Worker(signalChan chan Job, quit chan struct{}) {
	var lock sync.RWMutex

	utils.Log.Info("Starting worker routine")

	cacheChan := cache.NewCache(quit)

	periodicTimer := time.NewTicker(interval * time.Second)

	for {
		select {
		case <-quit:
			utils.Log.Info("Exiting Worker Routine")
			return
		case <-periodicTimer.C:
			utils.Log.Info("Worker alive")
		case request := <-signalChan:
			request.Task.Lock = &lock
			request.Task.CacheChan = cacheChan
			response, err := request.Task.RunHandler()

			request.Response <- JobResponse{Data: response, Err: err}
		}
	}
}
