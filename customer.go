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
	var s *string

	firstNameKey := genCustomerKey(custId, "firstName")
	s, err = Get(firstNameKey)
	checkAndLogError(err, firstNameKey)
	if err == nil && s != nil {
		c.firstName = *s
	}

	lastNameKey := genCustomerKey(custId, "lastName")
	s, err = Get(lastNameKey)
	checkAndLogError(err, lastNameKey)
	if err == nil && s != nil {
		c.lastName = *s
	}

	streetAddrKey := genCustomerKey(custId, "streetAddress")
	s, err = Get(streetAddrKey)
	checkAndLogError(err, streetAddrKey)
	if err == nil && s != nil {
		c.streetAddress = *s
	}

	cityKey := genCustomerKey(custId, "city")
	s, err = Get(cityKey)
	checkAndLogError(err, cityKey)
	if err == nil && s != nil {
		c.city = *s
	}

	stateKey := genCustomerKey(custId, "state")
	s, err = Get(stateKey)
	checkAndLogError(err, stateKey)
	if err == nil && s != nil {
		c.state = *s
	}

	zipKey := genCustomerKey(custId, "zip")
	s, err = Get(zipKey)
	checkAndLogError(err, zipKey)
	if err == nil && s != nil {
		c.zip = *s
	}

	return &c
}
