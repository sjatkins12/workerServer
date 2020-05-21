package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sjatkins12/workerServer/utils"
	"github.com/sjatkins12/workerServer/worker"
)

// SetupRouter ... register endpoint handlers with router
// send tasks to worker queue to be processed
func SetupRouter(queue chan worker.Job) *gin.Engine {
	router := gin.New()

	router.Use(utils.Ginrus(utils.Log), gin.Recovery())

	// addTaskToQueue ... Send a task to the worker routine
	addTaskToQueue := func(task string, data interface{}) chan worker.JobResponse {
		response := make(chan worker.JobResponse, 1)
		queue <- worker.Job{
			Response: response,
			Task:     worker.NewTask(task, data),
		}
		return response
	}

	router.GET("/howru", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/sampleTask", func(c *gin.Context) {
		responseChan := addTaskToQueue("sampleTask", "Lets show off")
		response := <-responseChan

		c.JSON(setHttpStatus(response.Err), gin.H{"message": response.Data})
	})

	return router
}

func setHttpStatus(err error) int {
	switch err {
	case nil:
		return http.StatusOK
	default:
		return http.StatusInternalServerError
	}
}
