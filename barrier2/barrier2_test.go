package barrier2

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
	t.Log("Test1")
	TEST := 10
	bc := make(chan bool, TEST)
	g := NewGroup(TEST)
	for i := 0; i < TEST; i++ {  
		go worker(i, g, t, bc)
	}
	
	for i := 0; i < TEST; i++ {
		<- bc
	}
}
const SIZE=10
var Test3data [SIZE] int

func worker2(n int, g * Group, t *testing.T) {
	Test3data[n] = n
	g.Wait()
}

func Test3(t *testing.T) {
	t.Log("Test3")
	GroupSize := SIZE + 1
	g := NewGroup(GroupSize) // +1 for me, who will wait for SIZE to be done

	// check initialization
	for i := 0; i < SIZE; i++ {
		if Test3data[i] != 0 {
			t.Errorf("Failed start condition on test3 (testdata not 0'd)\n")
		}
	}

	for i := 0; i < SIZE; i++ {
		go worker2(i, g, t)
	}
	g.Wait()  // if you comment this out, none of the array updates will be observed per the Go memory model

	t.Logf("Test3data = %v\n", Test3data)
	// check results
	for i := 0; i < SIZE; i++ {
		if Test3data[i] != i {
			t.Errorf("Test3data[%d] = %d, should be %d\n", i, Test3data[i], i)
		}
	}
}