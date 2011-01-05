package barrier

import (
	"sync"
)

type Group struct {
	count int // number of checkins
	total int
	signalin chan bool
	resize chan int
	waiter * sync.RWMutex
}

func controller (g * Group) {
	var r int
	for g.count != g.total {
		select {
		case <- g.signalin:
			g.count++
		case r = <- g.resize:
			g.total += r 
		}
	}

	g.waiter.Unlock()
}

func NewGroup (size int) (* Group) {
	g := &Group{0,size, make(chan bool, size), make(chan int), &sync.RWMutex{}}
	g.waiter.Lock()
	go controller(g)
	return g
}

func (g * Group) AddN (n int) {
	g.resize <- n
}

func (g * Group) Wait () {
	g.signalin <- true;
	g.waiter.RLock()
}

		

	
	
