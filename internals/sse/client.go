package sse

import (
	"file-analyzer/internals/domain"
	"sync"
)

type SSEManager struct {
	mu         sync.Mutex
	sseClients map[int64]chan domain.DocEvent
}

func NewSSEManager() *SSEManager {
	sseClientsMap := make(map[int64]chan domain.DocEvent)
	return &SSEManager{sseClients: sseClientsMap}
}

func (sse *SSEManager) Notify(event domain.DocEvent) {
	sse.mu.Lock()
	ch, ok := sse.sseClients[event.UserID]
	if !ok {
		ch = make(chan domain.DocEvent)
		sse.sseClients[event.UserID] = ch
	}
	sse.mu.Unlock()
	select {
	case ch <- event:
	default:
	}
}

func (sse *SSEManager) AddClient(key int64) chan domain.DocEvent {
	ch := make(chan domain.DocEvent, 10)
	sse.mu.Lock()
	sse.sseClients[key] = ch
	sse.mu.Unlock()

	return ch
}

func (sse *SSEManager) RemoveClient(key int64) {
	sse.mu.Lock()
	delete(sse.sseClients, key)
	sse.mu.Unlock()
}
