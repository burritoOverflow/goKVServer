package main

import "testing"

const custId int64 = 123456

var keyAttrs = [6]string{"firstName", "lastName", "streetAddress", "city", "state", "zip"}

// clean up after adding customer
func cleanUpAfterPutCustomer(t *testing.T) {
	for _, k := range keyAttrs {
		custKey := genCustomerKey(custId, k)
		err := Delete(custKey)
		if err != nil {
			t.Errorf("Got error deleting key: %s", custKey)
		}
	}
}

func TestGenCustomerKey(t *testing.T) {
	c := Customer{123456, "TestFirst", "TestLast", "123 Main Street", "Nowhere", "DE", "000000"}
	expectKey := "cust:123456:firstName"
	genKey := genCustomerKey(c.custId, "firstName")

	if genKey != expectKey {
		t.Errorf("Expected %s, got %s", expectKey, genKey)
	}
}

func TestPutGenKeyVal(t *testing.T) {
	cust := Customer{custId, "TestFirst", "TestLast", "123 Main Street", "Nowhere", "MD", "777777"}
	cust.AddCustomerRecord()
	defer cleanUpAfterPutCustomer(t)

	// store the results
	resMap := make(map[string]string)

	// ensure all K:V exist for this customer
	for _, k := range keyAttrs {
		custKey := genCustomerKey(custId, k)
		val, err := Get(custKey)

		if err != nil {
			t.Errorf("Error Getting key: %s", custKey)
		}
		// add val to map
		resMap[k] = val
	}

	if resMap["firstName"] != cust.firstName {
		t.Error("firstName incorrectly set")
	}

	if resMap["lastName"] != cust.lastName {
		t.Error("lastName incorrectly set")
	}

	if resMap["streetAddress"] != cust.streetAddress {
		t.Error("streetAdress incorrectly set")
	}

	if resMap["city"] != cust.city {
		t.Error("city incorrectly set")
	}

	if resMap["state"] != cust.state {
		t.Error("state incorrectly set")
	}

	if resMap["zip"] != cust.zip {
		t.Error("zip incorrectly set")
	}

	custTest := GetCustomerRecord(custId)
	if cust != *custTest {
		t.Errorf("Expect customer %s, got customer %s", cust.asStr(), custTest.asStr())
	}
}
