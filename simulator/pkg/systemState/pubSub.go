package systemState

import (
	"CPU-Simulator/simulator/pkg/logger"
	"sync"
)

type PubSub[T any] struct {
	Subscribers []chan T
	Mu          sync.RWMutex
}

func NewPubSub[T any]() *PubSub[T] {
	pubSub := &PubSub[T]{
		Subscribers: make([]chan T, 0),
	}
	return pubSub
}

func (pubSub *PubSub[T]) Subscribe() <-chan T {
	pubSub.Mu.Lock()
	defer pubSub.Mu.Unlock()

	channel := make(chan T)
	pubSub.Subscribers = append(pubSub.Subscribers, channel)

	return channel
}

func (pubSub *PubSub[T]) Publish(value T) {
	logger.Log.Println("INFO: Publish()")
	pubSub.Mu.RLock()
	defer pubSub.Mu.RUnlock()

	for i, channel := range pubSub.Subscribers {
		logger.Log.Println("INFO: Publish()", i)
		channel <- value
	}
}
