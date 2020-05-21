package cache

import (
	"sync"
	"time"

	"github.com/sjatkins12/workerServer/utils"
)

const INTERVAL time.Duration = 120

type CacheRequest struct {
	Type     string
	DataType string
	Data     interface{}
	Response chan CacheResponse
}

type CacheResponse struct {
	Data interface{}
	OK   bool
}

func NewCache(quit chan struct{}) chan CacheRequest {
	cacheChan := make(chan CacheRequest, 20)
	utils.Log.Info("Creating cache routine")
	go cache(cacheChan, quit)
	return cacheChan
}

func cache(requestChan chan CacheRequest, quit chan struct{}) {
	var lock sync.RWMutex
	periodicTimer := time.NewTicker(INTERVAL * time.Second)
	store := make(map[string]interface{})

	for {
		select {
		case <-quit:
			utils.Log.Info("Shutting down cache routine")
			return
		case request := <-requestChan:
			switch request.Type {
			case "DELETE":
				utils.Log.Info("Cache delete request received")
				lock.Lock()
				store = make(map[string]interface{})
				lock.Unlock()
			case "GET":
				lock.RLock()
				get(request, store)
				lock.RUnlock()
			case "POST":
				lock.Lock()
				set(request, store)
				lock.Unlock()
			}
		case <-periodicTimer.C:
			lock.Lock()
			store = make(map[string]interface{})
			utils.Log.Info("Deleted cache")
			lock.Unlock()
		}
	}
}

func get(request CacheRequest, store map[string]interface{}) {
	data, ok := store[request.DataType]

	request.Response <- CacheResponse{Data: data, OK: ok}
}

func set(request CacheRequest, store map[string]interface{}) {
	store[request.DataType] = request.Data
}
