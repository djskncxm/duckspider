package core

import (
	"github.com/djskncxm/duckspider/httpio"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
)

type Scheduler struct {
	RequestQueue *linkedlistqueue.Queue
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		RequestQueue: linkedlistqueue.New(),
	}
}

func (s *Scheduler) NextRequest() (*httpio.Request, bool) {
	if s.RequestQueue.Empty() {
		return nil, false
	}
	value, ok := s.RequestQueue.Dequeue()
	if !ok {
		return nil, false
	}
	req, ok := value.(*httpio.Request)
	if !ok {
		return nil, false
	}

	return req, true
}

func (s *Scheduler) EnRequest(request *httpio.Request) {
	s.RequestQueue.Enqueue(request)
}

func (s *Scheduler) IsEmpty() bool {
	return s.RequestQueue.Empty()
}
