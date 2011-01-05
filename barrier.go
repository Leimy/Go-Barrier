package barrier

import (
	"sync"
)

type Group struct {
	count int // number of checkins
	total int
	signalin chan bool
	resize chan int
	die chan bool
	waiter * sync.RWMutex
}

func controller (g * Group) {
	var r int
	for g.count != g.total {
		select {
		case <- g.die:
			break
		case <- g.signalin:
			g.count++
		case r = <- g.resize:
			g.total += r 
		}
	}

	g.waiter.Unlock()
}

func NewGroup (size int) (* Group) {
	g := &Group{0,size, make(chan bool, size), make(chan int), make(chan bool), &sync.RWMutex{}}
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

func (g * Group) Destroy () {
	if g.count != g.total {
		g.die <- true
	}
}

func (g * Group) Reset () {
	g.Destroy()
	g.count = 0
	g.waiter = &sync.RWMutex{}
	g.waiter.Lock()
	go controller(g)
}
		

	
	
