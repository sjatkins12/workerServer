package worker

import (
	"fmt"
	"sync"

	"github.com/sjatkins12/workerServer/cache"
	"github.com/sjatkins12/workerServer/utils"
)

// TaskType ... A struct representing a task to be processed by a worker
type TaskType struct {
	Name      string
	Data      interface{}
	Lock      *sync.RWMutex
	CacheChan chan cache.CacheRequest
}

// NewTask ... Return a new Task with a correct in memory mapping
func NewTask(name string, payload interface{}) TaskType {
	return TaskType{
		Name: name,
		Data: payload,
	}
}

var requestHandler = map[string](func(*TaskType) (interface{}, error)){
	"sampleTask": SampleTask,
}

// RunHandler ... Call the handler for the type of task
func (task *TaskType) RunHandler() (interface{}, error) {
	var response interface{}
	var err error

	response, err = requestHandler[task.Name](task)

	utils.Log.Info(fmt.Sprintf("Proccessing of task %s completed", task.Name))

	if err != nil {
		utils.Log.Error(fmt.Sprintf("Failed to process request: %s: ", task.Name), err)
		return nil, err
	}

	return response, nil
}
