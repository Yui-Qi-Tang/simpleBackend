package consistency

import (
	"fmt"
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
		go func() {
			sharedVariable++ // race condition here
		}()
	}
	runtime.Gosched() // yeild task complete in multiple cpu env.
	// sometime sharedVariable and expectVar are equal somtime are not.
	t.Logf("TestSharedDataRunInGorontine => sharedVariable :%d, expectedVar: %d", sharedVariable, expectedVar)
}

func TestSharedDataConsistencyWithSyncMutex(t *testing.T) {
	var sharedVariable int
	var expectedVar = 10000
	var mutex sync.Mutex
	c := make(chan bool)
	// Hope sharedVariabe is equal to expectedVar after goroutine add + 1 on sharedVariable in expectedVar times
	for i := 0; i < expectedVar; i++ {
		go func(gNum int) {
			mutex.Lock()
			sharedVariable++
			if sharedVariable == expectedVar {
				fmt.Println("GID: ", gNum, " Hit!")
			}
			mutex.Unlock()
			c <- true
		}(i)
	}

	// wait all of goroutine return
	for i := 0; i < expectedVar; i++ {
		<-c
	}

	defer func() {
		if sharedVariable == expectedVar {
			t.Log("TestSharedDataConsistencyWithSyncMutex PASS")
		} else {
			t.Error("TestSharedDataConsistencyWithSyncMutex FAIL", sharedVariable, ", ", expectedVar)
		}
	}()

}
