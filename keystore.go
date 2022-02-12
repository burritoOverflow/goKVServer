package main

import (
	"errors"
	"log"
)

var keyStore = make(map[string]string)

var ErrorNoSuchKey = errors.New("no such key")
var ErrorKeyExists = errors.New("existing key")

func Delete(key string) error {
	log.Printf("Delete: Request to delete for Key: %s\n", key)
	delete(keyStore, key)
	return nil
}

func Get(key string) (string, error) {
	log.Printf("Get: Request to get key %s\n", key)
	value, ok := keyStore[key]
	if !ok {
		log.Printf("Get: No Key %s found\n", key)
		return "", ErrorNoSuchKey
	}
	return value, nil
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
	log.Printf("Put: Request to get key %s\n", key)

	_, contains := keyStore[key]
	if contains {
		log.Printf("Key: %s already exits; not adding", key)
		return ErrorKeyExists
	}

	// otherwise, add the key
	keyStore[key] = value
	return nil
}
