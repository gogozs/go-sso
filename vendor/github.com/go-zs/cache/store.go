package cache

import (
	"sync"
)

const (
	defaultLength = 10000
)

type Store struct {
	l         *queue
	m         sync.Map
	maxLength int
}

type option struct {
	length int
}

type storeOption interface {
	apply(*option)
}

type optionFunc func(*option)

func (f optionFunc) apply(o *option) {
	f(o)
}

func WithLength(length int) optionFunc {
	return func(o *option) {
		o.length = length
	}
}

func NewStore(options ...storeOption) *Store {
	initialOption := option{length: defaultLength}
	for _, o := range options {
		o.apply(&initialOption)
	}
	return &Store{
		l:         NewQueue(),
		m:         sync.Map{},
		maxLength: initialOption.length,
	}
}
