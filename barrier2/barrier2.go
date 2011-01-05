package barrier2


import (
	"sync"
)

type Group struct {
	n int
	l * sync.Mutex
	waiter * sync.RWMutex
}

func NewGroup(n int) (* Group) {
	if n <= 0 {
		panic ("Group must be >= 1")
	}
	waiter := &sync.RWMutex{}
	waiter.Lock()
	return &Group{n, &sync.Mutex{}, waiter}
}

func (g* Group) Wait () {
	g.l.Lock() 
	g.n--
	if g.n == 0 {
		g.waiter.Unlock()
	}
	g.l.Unlock()
	g.waiter.RLock()
}
