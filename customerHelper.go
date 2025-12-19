package main

import (
	"strings"

	"github.com/lib/pq"
)

// getCustomerById retrieves a customer by their ID
func getCustomerById(customerId int) (customer Customer, err error) {

	err = db.QueryRow(GET_CUSTOMER_BY_ID, customerId).Scan(
		&customer.ID,
		&customer.Address,
		&customer.CreatedAt,
		&customer.Info,
		&customer.Mobile,
		&customer.Name,
		&customer.ReferredBy, // Will be -1 if NULL
		&customer.UpdatedAt,
	)

	if err != nil {
		return customer, err
	}
	return customer, nil
}

// isForeignKeyViolation checks if an error is a foreign key constraint violation
func isForeignKeyViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23503" // foreign_key_violation
	}
	return strings.Contains(strings.ToLower(err.Error()), "foreign key")
}
