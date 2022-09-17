package main

import (
	"sort"
	"testing"
)

func TestKeyStorePut(t *testing.T) {
	InitKeyStore()
	keyStr := "keyone"
	valStr := "valone"
	defer delete(keyStore.m, keyStr)

	err := Put(keyStr, valStr)
	if err != nil {
		t.Error("Put failed")
	}

	_, contains := keyStore.m[keyStr]
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
	InitKeyStore()
	testKey := "testkey"
	testVal := "testval"

	defer delete(keyStore.m, testKey)
	keyStore.m[testKey] = testVal

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
	InitKeyStore()
	testKey := "testkey"
	testVal := "testval"

	keyStore.m[testKey] = testVal

	err := Delete(testKey)
	if err != nil {
		t.Error("Got error deleting key")
	}

	_, contains := keyStore.m[testKey]
	if contains {
		t.Error("Delete failed; map contains k:v")
	}
}

func TestUpdate(t *testing.T) {
	// check that an existing key's update works
	InitKeyStore()
	testKey := "testkey"
	testVal := "testval"

	keyStore.m[testKey] = testVal

	newVal := "newval"
	err := Update(testKey, newVal)
	if err != nil {
		t.Error("Error updating key")
	}

	if keyStore.m[testKey] != newVal {
		t.Error("Key not updated to new value")
	}

	err = Update("notvalidkey", "notvalidvalue")
	if err == nil {
		t.Error("Updating non-existant key should result in error")
	}
}

func TestGetAll(t *testing.T) {
	// remove all existing keys before testing
	InitKeyStore()
	keyStore.m = make(map[string]string)

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
	sort.Slice(kvList, func(i, j int) bool {
		return kvList[i].Key < kvList[j].Key
	})

	keyStore.m[testKeyOne] = testValOne
	keyStore.m[testKeyTwo] = testValTwo
	keyStore.m[testKeyThree] = testValThree

	results := GetAll()
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	ok := testEq(kvList, results)
	if !ok {
		t.Errorf("Contents not equal, expected: %s, got: %s", kvList, results)
	}
}

func TestDelete(t *testing.T) {
	InitKeyStore()
	key := "key"
	value := "value"
	keyStore.m[key] = value

	err := Delete(key)
	if err != nil {
		t.Error("Error when attempting delete")
	}

	// ensure that the key was removed
	_, contains := keyStore.m[key]
	if contains {
		t.Errorf("Key %s not deleted after Delete\n", key)
	}

	// attempts to delete existing key should throw error
	err = Delete(key)
	if err == nil {
		t.Error("No error thrown when Delete on non-existing key")
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
