package main

import (
	"errors"
	"log"
)

var store = make(map[string]string)

var ErrorNoSuchKey = errors.New("no such key")

func Delete(key string) error {
	log.Printf("Delete: Request to delete for Key: %s\n", key)
	delete(store, key)
	return nil
}

func Get(key string) (string, error) {
	log.Printf("Get: Request to get key %s\n", key)
	value, ok := store[key]
	if !ok {
		log.Printf("Get: No Key %s found\n", key)
		return "", ErrorNoSuchKey
	}
	return value, nil
}

func GetAll() KVList {
	var kvs KVList
	for k, v := range store {
		kvs = append(kvs, KeyValEntry{Key: k, Value: v})
	}
	return kvs
}

func Put(key string, value string) error {
	log.Printf("Put: Request to get key %s\n", key)
	store[key] = value
	return nil
}
