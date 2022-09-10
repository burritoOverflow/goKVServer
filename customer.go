package main

import (
	"fmt"
	"log"
)

type Customer struct {
	custId        int64
	firstName     string
	lastName      string
	streetAddress string
	city          string
	state         string
	zip           string
}

func checkAndLogError(err error, keyname string) {
	if err != nil {
		log.Printf("Error: %s - Unable to retrieve for key: %s", err, keyname)
	}
}

func (c *Customer) asStr() string {
	return fmt.Sprintf("%d, %s, %s, %s, %s, %s, %s", c.custId, c.firstName, c.lastName, c.streetAddress, c.city, c.state, c.zip)
}

// Generate the customer key
// where K = `cust:CUST_ID:_CUST_FIELD_NAME`
func genCustomerKey(id int64, fieldName string) string {
	strtKey := "cust"
	return fmt.Sprintf("%s:%d:%s", strtKey, id, fieldName)
}

// Add K:V pairs for each field of the Customer provided
// such that K = `cust:CUST_ID:_CUST_FIELD_NAME`, eg `cust:123:firstName`
// for each field
func (c *Customer) AddCustomerRecord() (err error) {
	err = Put(genCustomerKey(c.custId, "firstName"), c.firstName)
	err = Put(genCustomerKey(c.custId, "lastName"), c.lastName)
	err = Put(genCustomerKey(c.custId, "streetAddress"), c.streetAddress)
	err = Put(genCustomerKey(c.custId, "city"), c.city)
	err = Put(genCustomerKey(c.custId, "state"), c.state)
	err = Put(genCustomerKey(c.custId, "zip"), c.zip)
	return
}

func GetCustomerRecord(custId int64) *Customer {
	c := Customer{custId: custId}
	var err error

	firstNameKey := genCustomerKey(custId, "firstName")
	c.firstName, err = Get(firstNameKey)
	checkAndLogError(err, firstNameKey)

	lastNameKey := genCustomerKey(custId, "lastName")
	c.lastName, err = Get(lastNameKey)
	checkAndLogError(err, lastNameKey)

	streetAddrKey := genCustomerKey(custId, "streetAddress")
	c.streetAddress, err = Get(streetAddrKey)
	checkAndLogError(err, streetAddrKey)

	cityKey := genCustomerKey(custId, "city")
	c.city, err = Get(cityKey)
	checkAndLogError(err, cityKey)

	stateKey := genCustomerKey(custId, "state")
	c.state, err = Get(stateKey)
	checkAndLogError(err, stateKey)

	zipKey := genCustomerKey(custId, "zip")
	c.zip, err = Get(zipKey)
	checkAndLogError(err, zipKey)

	return &c
}
