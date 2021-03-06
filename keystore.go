package main

import (
	"errors"
	"log"
)

var keyStore = make(map[string]string)

var ErrorNoSuchKey = errors.New("no such key")
var ErrorKeyExists = errors.New("existing key")

// Delete the key from the map; err if not found
func Delete(key string) error {
	log.Printf("Delete: Request to delete for Key: %s\n", key)
	// delete doesn't return err, but inform the user of a bad req
	_, contains := keyStore[key]
	if !contains {
		log.Printf("Delete: Cannot delete non-existant key %s\n", key)
		return ErrorNoSuchKey
	}
	delete(keyStore, key)
	return nil
}

// Return the key from the map if found, err otherwise
func Get(key string) (string, error) {
	log.Printf("Get: Request to get key %s\n", key)
	value, ok := keyStore[key]
	if !ok {
		log.Printf("Get: No Key %s found\n", key)
		return "", ErrorNoSuchKey
	}
	return value, nil
}

// Update key to value, only if key exists
func Update(key string, value string) error {
	_, contains := keyStore[key]
	if !contains {
		// key doesn't exist, cannot update
		log.Printf("Update: Error updating key %s, does not exist in store\n", key)
		return ErrorNoSuchKey
	}
	keyStore[key] = value
	return nil
}

func GetAll() KVList {
	var kvs KVList
	for k, v := range keyStore {
		kvs = append(kvs, KeyValEntry{Key: k, Value: v})
	}
	return kvs
}

// Only allow put to succeed when the key does not exist
func Put(key string, value string) error {
	log.Printf("Put: Request to put key %s\n", key)

	_, contains := keyStore[key]
	if contains {
		log.Printf("Put: Key: %s already exits; not adding", key)
		return ErrorKeyExists
	}

	// otherwise, add the key
	keyStore[key] = value
	return nil
}
