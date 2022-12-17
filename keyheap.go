package main

import (
	"container/heap"
	"errors"
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

func (kmh *KeyMinHeap) Delete(key string) (err error) {
	if kmh.Len() == 0 {
		err = errors.New("KeyMinHeap has length 0; not instantiated prior to Delete")
		return
	}

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
	return
}

func (kmh KeyMinHeap) Less(i, j int) bool {
	return kmh[i].timestamp.Before(kmh[j].timestamp)
}

func (kmh KeyMinHeap) Swap(i, j int) {
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
