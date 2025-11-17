package utils

type Semaphore struct {
	queue chan int
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		queue: make(chan int, n),
	}
}

func (s *Semaphore) Acquire() {
	s.queue <- 1
}

func (s *Semaphore) Release() {
	<-s.queue
}
