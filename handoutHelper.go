package main

import "errors"

func validateHandout(handout HandoutUpdate) error {

	if handout.UserId == 0 {
		return errors.New("user cannot be empty")
	}

	if handout.Date.IsZero() {
		return errors.New("date cannot be empty")
	}

	if handout.Amount <= 0 {
		return errors.New("enter a valid amount")
	}
	return nil
}
