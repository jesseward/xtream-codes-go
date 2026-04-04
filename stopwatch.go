package xtream_codes_go

import (
	"sync"
	"time"
)

type stopwatch struct {
	mu     sync.RWMutex
	events map[string]*stopwatchEvent
}

func (s *stopwatch) getContext() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var ctx = make(map[string]interface{})

	for name, event := range s.events {
		ctx[name] = event.String()
	}

	return ctx
}

func (s *stopwatch) GetEvent(name string) *stopwatchEvent {
	s.mu.RLock()
	event, ok := s.events[name]
	s.mu.RUnlock()

	if !ok {
		s.mu.Lock()
		// Double check after acquiring write lock
		event, ok = s.events[name]
		if !ok {
			event = new(stopwatchEvent)
			s.events[name] = event
		}
		s.mu.Unlock()
	}

	return event
}

type stopwatchEvent struct {
	start time.Time
	stop  time.Time
}

func (s *stopwatchEvent) Start() {
	s.start = time.Now()
}

func (s *stopwatchEvent) Stop() {
	s.stop = time.Now()
}

func (s *stopwatchEvent) String() string {
	return s.stop.Sub(s.start).String()
}
