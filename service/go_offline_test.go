package service_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoOffline", func() {
 	It("should successfully set the driver's status to offline", func() {
 		// Implement success case
 		// 1. Create a mock driver with status online
 		// 2. Call the go_offline function with the mock driver
 		// 3. Check that the driver's status is now offline
 	})
 
 	It("should return an error when the driver does not exist", func() {
 		// Implement failure case
 		// 1. Call the go_offline function with a non-existent driver
 		// 2. Check that the function returns an error
 	})
})

