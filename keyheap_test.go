package main

import (
	"container/heap"
	"fmt"
	"testing"
	"time"
)

func TestAtKeyMinHeap(t *testing.T) {
	keyMinHeap := &KeyMinHeap{}
	heap.Init(keyMinHeap)

	n := 100
	for i := 0; i < n; i++ {
		heap.Push(keyMinHeap, fmt.Sprintf("Key:%d", i))
	}

	for i := 0; i < n; i++ {
		keyDate := keyMinHeap.At(i)
		kdStr := fmt.Sprintf("Key:%d", i)

		if kdStr != keyDate.Key {
			t.Errorf("Incorrect key, expected: %s, got %s", kdStr, keyDate.Key)
		}
	}
}

func TestSwapKeyMinHeap(t *testing.T) {
	keyMinHeap := &KeyMinHeap{}
	heap.Init(keyMinHeap)
	keyOneStr := "KeyOne"
	keyTwoStr := "KeyTwo"

	keyMinHeap.Push(keyOneStr)
	keyMinHeap.Push(keyTwoStr)

	if keyMinHeap.At(0).Key != keyOneStr && keyMinHeap.At(1).Key != keyTwoStr {
		t.Errorf("Insertion incorrect")
	}

	keyMinHeap.Swap(0, 1)
	if keyMinHeap.At(0).Key != keyTwoStr && keyMinHeap.At(1).Key != keyOneStr {
		t.Errorf("Swap results incorrect")
	}
}

func TestKeyMinHeapIdxOf(t *testing.T) {
	keyMinHeap := &KeyMinHeap{}
	heap.Init(keyMinHeap)

	n := 100
	for i := 0; i < n; i++ {
		heap.Push(keyMinHeap, fmt.Sprintf("Key:%d", i))
	}

	for i := 0; i < n; i++ {
		res := keyMinHeap.idxOf(fmt.Sprintf("Key:%d", i))
		if res != i {
			t.Errorf("Expected %d as idxOf, got %d", i, res)
		}
	}
}

func TestLessKeyMinHeap(t *testing.T) {
	keyMinHeap := &KeyMinHeap{}
	heap.Init(keyMinHeap)
	n := 35

	for i := 0; i < n; i++ {
		time.Sleep(10 * time.Millisecond)
		// add in order; sleep before each to have a testable ordered slice
		heap.Push(keyMinHeap, fmt.Sprintf("Key:%d", i))
	}

	for i := 0; i < n-1; i++ {
		// ensure each are ordered correctly
		if keyMinHeap.Less(i, i+1) != true {
			t.Errorf("Incorrect ordering")
		}
	}
}

func TestKeyMinHeapPop(t *testing.T) {
	// must instantiate as heap
	keyMinHeap := &KeyMinHeap{}
	heap.Init(keyMinHeap)
	n := 12

	for i := 0; i < n; i++ {
		heap.Push(keyMinHeap, fmt.Sprintf("Key:%d", i))
		time.Sleep(10 * time.Millisecond)
	}

	var prevTime time.Time
	for i := 0; i < n; i++ {
		expectedStr := fmt.Sprintf("Key:%d", i)
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

func TestKeyHeapKeyStoreLimit(t *testing.T) {
	InitKeyStore()

	heapLimit := 12
	for i := 0; i < heapLimit; i++ {
		_ = Put(fmt.Sprintf("Key:%d", i), "val")
	}

	for i := 0; i < heapLimit; i++ {
		// add new keys past the limit
		_ = Put(fmt.Sprintf("Key:%d", 100+i), "val")
		if keyStore.kmh.Len() != heapLimit {
			t.Errorf("Expected heap limit to be %d, got Len %d", heapLimit, keyStore.kmh.Len())
		}

		// check that the original keys are no longer present
		oldKey := fmt.Sprintf("Key:%d", i)
		_, err := Get(oldKey)
		if err == nil {
			t.Errorf("Got key %s - expected to be popped", oldKey)
		}
	}
}
