package main

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
