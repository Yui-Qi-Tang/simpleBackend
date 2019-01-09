package consistency

import (
	"runtime"
	"sync"
	"testing"
	_ "testing"
	_ "time"
)

func TestSharedDataRunInGorontine(t *testing.T) {
	var sharedVariable int
	var expectedVar = 100
	// Hope sharedVariabe is equal to expectedVar after goroutine add + 1 on sharedVariable in expectedVar times, but...
	for i := 0; i < expectedVar; i++ {
		go func(x int) {
			sharedVariable++ // race condition here
		}(i)
	}
	runtime.Gosched() // yeild task complete in multiple cpu env.
	// sometime sharedVariable and expectVar are equal somtime are not.
	t.Logf("TestSharedDataRunInGorontine => sharedVariable :%d, expectedVar: %d", sharedVariable, expectedVar)
}

func TestSharedDataConsistencyWithSyncWG(t *testing.T) {
	var sharedVariable int
	var expectedVar = 100
	var wg sync.WaitGroup // sync wait group

	// Hope sharedVariabe is equal to expectedVar after goroutine add + 1 on sharedVariable in expectedVar times
	for i := 0; i < expectedVar; i++ {
		wg.Add(1)
		go func(x int) {
			sharedVariable++
			wg.Done()
		}(i)
	}
	wg.Wait()
	if sharedVariable == expectedVar {
		t.Log("TestSharedDataConsistencyWithSyncWG PASS")
	} else {
		t.Error("TestSharedDataConsistencyWithSyncWG FAIL", sharedVariable, ", ", expectedVar)
	}
}
