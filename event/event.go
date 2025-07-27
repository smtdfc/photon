package event

import (
	"sync"
	"sync/atomic"
)

type EventData map[string]any
type EventListener func(data EventData)

type listener struct {
	id   int64
	fn   EventListener
}

type EventEmitter struct {
	events     map[string][]listener
	lock       sync.RWMutex
	nextID     int64
}

func CreateEventEmitter() *EventEmitter {
	return &EventEmitter{
		events: make(map[string][]listener),
	}
}

func (e *EventEmitter) AddListener(name string, callback EventListener) int64 {
	id := atomic.AddInt64(&e.nextID, 1)

	e.lock.Lock()
	defer e.lock.Unlock()

	e.events[name] = append(e.events[name], listener{id: id, fn: callback})
	return id
}

func (e *EventEmitter) RemoveListener(name string, id int64) {
	e.lock.Lock()
	defer e.lock.Unlock()

	listeners := e.events[name]
	for i, l := range listeners {
		if l.id == id {
			e.events[name] = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}
}

func (e *EventEmitter) Emit(name string, data EventData) {
	e.lock.RLock()
	listeners := append([]listener{}, e.events[name]...)
	e.lock.RUnlock()

	var wg sync.WaitGroup
	for _, l := range listeners {
		wg.Add(1)
		fn := l.fn
		go func() {
			defer wg.Done()
			fn(data)
		}()
	}
	wg.Wait()
}