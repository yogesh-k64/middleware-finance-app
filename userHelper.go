package main

import (
	"strings"

	"github.com/lib/pq"
)

func getUserById(userId int) (user User, err error) {

	err = db.QueryRow(GET_USER_BY_ID, userId).Scan(
		&user.ID,
		&user.Address,
		&user.CreatedAt,
		&user.Info,
		&user.Mobile,
		&user.Name,
		&user.ReferredBy, // Will be -1 if NULL
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}
	return user, nil
}

func isForeignKeyViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23503" // foreign_key_violation
	}
	return strings.Contains(strings.ToLower(err.Error()), "foreign key")
}
