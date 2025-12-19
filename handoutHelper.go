package main

import "errors"

func validateHandout(handout HandoutUpdate) error {

	if handout.CustomerId == 0 {
		return errors.New("customer cannot be empty")
	}

	if handout.Date.IsZero() {
		return errors.New("date cannot be empty")
	}

	if handout.Amount <= 0 {
		return errors.New("enter a valid amount")
	}
	return nil
}

func validateCollection(collection Collection) error {

	if collection.HandoutId == 0 {
		return errors.New("handout id cannot be empty")
	}

	if collection.Date.IsZero() {
		return errors.New("date cannot be empty")
	}

	if collection.Amount <= 0 {
		return errors.New("enter a valid amount")
	}
	return nil
}
