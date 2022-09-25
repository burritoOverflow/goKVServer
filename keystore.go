package main

import (
	"container/heap"
	"errors"
	"log"
	"sync"
)

const maxKeys int = 12

var keyStore = struct {
	m   map[string]string
	kmh KeyMinHeap
	sync.RWMutex
}{m: make(map[string]string), kmh: KeyMinHeap{}}

var ErrorNoSuchKey = errors.New("no such key")
var ErrorKeyExists = errors.New("existing key")

func popKeyHeap() string {
	popVal := heap.Pop(&keyStore.kmh)
	return popVal.(KeyDate).Key
}

func pushKeyHeap(key string) {
	heap.Push(&keyStore.kmh, key)
}

// InitKeyStore create the heap associated with the keystore
func InitKeyStore() {
	keyStore.m = make(map[string]string)
	heap.Init(&keyStore.kmh)
}

func lenKeyStore() int {
	return len(keyStore.m)
}

// Delete the key from the map; err if not found
func Delete(key string) error {
	log.Printf("Delete: Request to delete for Key: %s\n", key)
	// delete doesn't return err, but inform the user of a bad req
	keyStore.Lock()
	_, contains := keyStore.m[key]
	if !contains {
		keyStore.Unlock()
		log.Printf("Delete: Cannot delete non-existant key %s\n", key)
		return ErrorNoSuchKey
	}
	delete(keyStore.m, key)
	keyStore.kmh.Delete(key)
	log.Printf("Deleted key %s", key)
	delete(keyStore.m, key)
	keyStore.Unlock()

	return nil
}

// Get Return the key from the map if found, err otherwise
func Get(key string) (*string, error) {
	log.Printf("Get: Request to get key %s\n", key)
	keyStore.RLock()
	value, ok := keyStore.m[key]
	keyStore.RUnlock()

	if !ok {
		log.Printf("Get: No Key %s found\n", key)
		return nil, ErrorNoSuchKey
	}
	return &value, nil
}

// Update key to value, only if key exists
func Update(key string, value string) error {
	keyStore.Lock()
	_, contains := keyStore.m[key]
	if !contains {
		// key doesn't exist, cannot update
		keyStore.Unlock()
		log.Printf("Update: Error updating key %s, does not exist in store\n", key)
		return ErrorNoSuchKey
	}

	keyStore.m[key] = value
	keyStore.Unlock()
	return nil
}

func GetAll() KVList {
	kvs := KVList{}
	keyStore.RLock()
	for k, v := range keyStore.m {
		kvs = append(kvs, KeyValEntry{Key: k, Value: v})
	}
	keyStore.RUnlock()
	return kvs
}

// Put Only allow put to succeed when the key does not exist
func Put(key string, value string) error {
	log.Printf("Put: Request to put key %s\n", key)
	keyStore.Lock()

	_, contains := keyStore.m[key]
	if contains {
		keyStore.Unlock()
		log.Printf("Put: Key: %s already exits; not adding", key)
		return ErrorKeyExists
	}

	// otherwise, add the key
	keyStore.m[key] = value
	if lenKeyStore() > maxKeys {
		popVal := popKeyHeap()
		log.Printf("Key Store reached limit; popped and removed %s", popVal)
		delete(keyStore.m, popVal)
	}
	pushKeyHeap(key)

	keyStore.m[key] = value
	keyStore.Unlock()
	return nil
}
