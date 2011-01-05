package barrier

import (
	"testing"
)

func worker (n int, g * Group, t *testing.T, bc chan <- bool) {
	t.Logf("[%d] <---\n", n)
	g.Wait()
	t.Logf("[%d] --->\n", n)
	bc <- true
}

func Test1(t *testing.T) {
	TEST := 100
	bc := make(chan bool)
	g := NewGroup(TEST)
	for i := 0; i < TEST; i++ {  
		go worker(i, g, t, bc)
	}
	
	for i := 0; i < TEST; i++ {
		<- bc
	}
}

func Test2(t *testing.T) {
	TEST_START := 100
	bc := make(chan bool)
	g := NewGroup(TEST_START)
	// add some in
	g.AddN(10)
	for i := 0; i < TEST_START+10; i++ {
		go worker(i, g, t, bc)
	}

	for i := 0; i < TEST_START+10; i++ {
		<- bc
	}
}
