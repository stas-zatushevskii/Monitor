package pool

import (
	"sync"
)

type Resetter interface {
	Reset()
}

type Pool[T Resetter] struct {
	p sync.Pool
}

func New[T Resetter](construct func() T) *Pool[T] {
	return &Pool[T]{
		p: sync.Pool{
			New: func() any { return construct() },
		},
	}
}

func (p *Pool[T]) Get(obj Resetter) Resetter {
	return p.Get(obj)
}

func (p *Pool[T]) Put(obj T) {
	obj.Reset()
	p.p.Put(obj)
}
