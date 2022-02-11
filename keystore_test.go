package main

import (
	"reflect"
	"testing"
)

func TestKeyStorePut(t *testing.T) {
	keyStr := "keyone"
	valStr := "valone"
	defer delete(keyStore, keyStr)

	err := Put(keyStr, valStr)
	if err != nil {
		t.Error("Put failed")
	}

	_, contains := keyStore[keyStr]
	if !contains {
		t.Errorf("Value %s not stored in keystore with key %s\n", valStr, keyStr)
	}
}

func TestKeyStoreGet(t *testing.T) {
	testKey := "testkey"
	testVal := "testval"

	defer delete(keyStore, testKey)
	keyStore[testKey] = testVal

	val, err := Get(testKey)
	if err != nil {
		t.Error("Error when getting key from keystore")
	}

	if val != testVal {
		t.Error("Value not as expected")
	}

	nonKey := "nonkey"
	_, err = Get(nonKey)
	if err == nil {
		t.Error("Error not returned when requesting non-existent key")
	}
}

func TestKeyStoreDelete(t *testing.T) {
	testKey := "testkey"
	testVal := "testval"

	keyStore[testKey] = testVal

	err := Delete(testKey)
	if err != nil {
		t.Error("Got error deleting key")
	}

	_, contains := keyStore[testKey]
	if contains {
		t.Error("Delete failed; map contains k:v")
	}
}

func TestGetAll(t *testing.T) {
	testKeyOne := "one"
	testValOne := "valone"
	testKeyTwo := "two"
	testValTwo := "valtwo"
	testKeyThree := "keythree"
	testValThree := "valthree"

	tkOneKV := KeyValEntry{Key: testKeyOne, Value: testValOne}
	tkTwoKv := KeyValEntry{Key: testKeyTwo, Value: testValTwo}
	tkThreeKv := KeyValEntry{Key: testKeyThree, Value: testValThree}

	kvList := KVList{tkOneKV, tkTwoKv, tkThreeKv}

	keyStore[testKeyOne] = testValOne
	keyStore[testKeyTwo] = testValTwo
	keyStore[testKeyThree] = testValThree

	results := GetAll()

	if !reflect.DeepEqual(results, kvList) {
		t.Error("Contents not equal to test set")
	}
}
