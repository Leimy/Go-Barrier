package barrier

import (
	"testing"
)

func worker(n int, b *Barrier, t *testing.T, bc chan<- bool) {
	t.Logf("[%d] <---\n", n)
	b.Wait()
	t.Logf("[%d] --->\n", n)
	bc <- true
}

func Test1(t *testing.T) {
	t.Log("Test1")
	TEST := 10
	bc := make(chan bool, TEST)
	g := NewBarrier(TEST)
	for i := 0; i < TEST; i++ {
		go worker(i, g, t, bc)
	}

	for i := 0; i < TEST; i++ {
		<-bc
	}
}

const SIZE = 10

var Test3data [SIZE]int

func worker2(n int, b *Barrier, t *testing.T) {
	Test3data[n] = n
	b.Wait()
}

func Test3(t *testing.T) {
	t.Log("Test3")
	GroupSize := SIZE + 1
	b := NewBarrier(GroupSize) // +1 for me, who will wait for SIZE to be done

	// check initialization
	for i := 0; i < SIZE; i++ {
		if Test3data[i] != 0 {
			t.Errorf("Failed start condition on test3 (testdata not 0'd)\n")
		}
	}

	for i := 0; i < SIZE; i++ {
		go worker2(i, b, t)
	}
	b.Wait() // if you comment this out, none of the array updates will be observed per the Go memory model

	t.Logf("Test3data = %v\n", Test3data)
	// check results
	for i := 0; i < SIZE; i++ {
		if Test3data[i] != i {
			t.Errorf("Test3data[%d] = %d, should be %d\n", i, Test3data[i], i)
		}
	}
}

func Test4(t *testing.T) {
	size := 10
	Test4data := make([]int, size)

	GroupSize := size
	b := NewBarrier(GroupSize)

	for i := 0; i < size-1; i++ {
		go func(i int) {
			t.Logf("I'm %d and I'm setting index %d to %d\n", i, i, i)
			Test4data[i] = i
			b.Wait()
		}(i)
	}
	Test4data[size-1] = size - 1
	t.Logf("Test4data: about to wait\n")
	b.Wait() // if you comment this out, none of the array updates will be observed per the Go memory model

	t.Logf("Test4data: %v\n", Test4data)
	// check results
	for i := 0; i < size; i++ {
		if Test4data[i] != i {
			t.Errorf("Test4data[%d] = %d, should be %d\n", i, Test4data[i], i)
		}
	}

	// Reset and run again, with different work
	b = NewBarrier(GroupSize)

	for i := 0; i < size-1; i++ {
		go func(i, size int) {
			t.Logf("index is: %d\n", i)
			Test4data[i] = size - i
			b.Wait()
		}(i, size)

	}
	Test4data[size-1] = 1
	b.Wait()

	t.Logf("Test4data: %v\n", Test4data)
	// check results
	for i := 0; i < size; i++ {
		t.Logf("Test4data[%d] = %d\n", i, Test4data[i])
		if Test4data[i] != size-i {
			t.Errorf("Test4data[%d] = %d, should be %d\n", i, Test4data[i], size-i)
		}
	}
}
