package main

import (
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

	// try put again; should fail with the existing key
	err = Put(keyStr, valStr)
	if err == nil {
		t.Errorf("Error should be present for existing key %s\n", keyStr)
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

func TestUpdate(t *testing.T) {
	// check that an existing key's update works
	testKey := "testkey"
	testVal := "testval"

	keyStore[testKey] = testVal

	newVal := "newval"
	err := Update(testKey, newVal)
	if err != nil {
		t.Error("Error updating key")
	}

	if keyStore[testKey] != newVal {
		t.Error("Key not updated to new value")
	}

	err = Update("notvalidkey", "notvalidvalue")
	if err == nil {
		t.Error("Updating non-existant key should result in error")
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

	ok := testEq(kvList, results)
	if !ok {
		t.Error("Contents not equal")
	}

}

func testEq(a, b KVList) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
