package main

import (
	"container/heap"
	"fmt"
	"testing"
	"time"
)

func TestKeyMinHeapPop(t *testing.T) {
	// must instantiate as heap
	keyMinHeap := &KeyMinHeap{}
	heap.Init(keyMinHeap)
	n := 12

	for i := 0; i < n; i++ {
		//keyMinHeap.Push(fmt.Sprintf("Key:%d", i))
		heap.Push(keyMinHeap, fmt.Sprintf("Key:%d", i))
		time.Sleep(10 * time.Millisecond)
	}

	var prevTime time.Time
	for i := 0; i < n; i++ {
		expectedStr := fmt.Sprintf("Key:%d", i)
		//popVal := keyMinHeap.Pop().(KeyDate)
		popVal := heap.Pop(keyMinHeap).(KeyDate)
		receivedStr := popVal.Key

		// check timestamp for last
		if i != 0 {
			// each time stamp should be
			if !prevTime.Before(popVal.timestamp) {
				t.Errorf("Key: %s Idx: %d - Popped timestamp %s did not occur after previous stored timestamp %s", receivedStr, i, popVal.timestamp, prevTime)
			}
		}
		prevTime = popVal.timestamp

		if receivedStr != expectedStr {
			t.Errorf("Expected %s, Got %s", expectedStr, receivedStr)
		}
	}
}
