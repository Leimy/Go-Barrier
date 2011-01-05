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

 func Test2(t *testing.T) {
	t.Log("Test2")
	TEST_START := 10
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

	// check results
	for i := 0; i < SIZE; i++ {
		t.Logf("Test3data[%d] = %d\n", i, Test3data[i])
		if Test3data[i] != i {
			t.Errorf("Test3data[%d] = %d, should be %d\n", i, Test3data[i], i)
		}
	}
}

func Test4(t *testing.T) {
	size := 10
	Test4data := make([]int, size)

	GroupSize := size + 1
	g := NewGroup(GroupSize)

	for i := 0; i < size; i++ {
	 	go func(i int) {
	 		Test4data[i] = i
	 		g.Wait()
	 	}(i)
	}
	g.Wait()  // if you comment this out, none of the array updates will be observed per the Go memory model

	t.Logf("Test4data: %v\n", Test4data)
	// check results
	for i := 0; i < size; i++ {
		if Test4data[i] != i {
			t.Errorf("Test4data[%d] = %d, should be %d\n", i, Test4data[i], i)
		}
	}

	// Reset and run again, with different work
	g.Reset()
	
//	g = NewGroup(GroupSize)
	for i := 0; i < size; i++ {
		go func (i, size int) {
			t.Logf("index is: %d\n", i)
			Test4data[i] = size-i
			g.Wait()
		}(i, size)
		
	}
	g.Wait()

	t.Logf("Test4data: %v\n", Test4data)
	// check results
	for i := 0; i < size; i++ {
		t.Logf("Test4data[%d] = %d\n", i, Test4data[i])
		if Test4data[i] != size - i {
			t.Errorf("Test4data[%d] = %d, should be %d\n", i, Test4data[i], size-i)
		}
	}
}