package main

import (
	"container/heap"
	"time"
)

type KeyDate struct {
	Key       string
	timestamp time.Time
	index     int
}

type KeyMinHeap []KeyDate

func (kmh KeyMinHeap) Len() int {
	return len(kmh)
}

func (kmh KeyMinHeap) idxOf(key string) int {
	for i, kd := range kmh {
		if kd.Key == key {
			return i
		}
	}
	return -1
}

func (kmh *KeyMinHeap) Delete(key string) {
	h := make(KeyMinHeap, kmh.Len()-1)
	j := 0

	for i := 0; i < kmh.Len()-1; i++ {
		keyMinHeap := *kmh
		if keyMinHeap[i].Key == key {
			continue
		}
		h[j] = keyMinHeap[i]
		j++
	}

	heap.Init(&h)
	*kmh = h
}

func (kmh KeyMinHeap) Less(i int, j int) bool {
	return kmh[i].timestamp.Before(kmh[j].timestamp)
}

func (kmh KeyMinHeap) Swap(i int, j int) {
	kmh[i], kmh[j] = kmh[j], kmh[i]
	kmh[i].index = i
	kmh[j].index = j
}

func (kmh *KeyMinHeap) Push(keyStr interface{}) {
	*kmh = append(*kmh, KeyDate{Key: keyStr.(string), timestamp: time.Now()})
}

func (kmh *KeyMinHeap) Pop() interface{} {
	old := *kmh
	n := len(old)
	item := old[n-1]
	item.index = -1
	*kmh = old[0 : n-1]
	return item
}
