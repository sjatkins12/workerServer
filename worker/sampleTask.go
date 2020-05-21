package worker

import "fmt"

func SampleTask(task *TaskType) (interface{}, error) {
	return fmt.Sprintf("%v : Now you make a task!", task.Data), nil
}
