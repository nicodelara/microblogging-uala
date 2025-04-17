package application

import "errors"

var (
	// ErrUserAlreadyExists is returned when trying to create a user that already exists
	ErrUserAlreadyExists = errors.New("username already exists")
	// ErrEmailAlreadyExists is returned when trying to create a user with an email that already exists
	ErrEmailAlreadyExists = errors.New("email already exists")
	// ErrFollowerNotFound is returned when the user trying to follow does not exist
	ErrFollowerNotFound = errors.New("follower user not found")
	// ErrFollowingNotFound is returned when the user to follow does not exist
	ErrFollowingNotFound = errors.New("user to follow not found")
	// ErrAlreadyFollowing is returned when already following the user
	ErrAlreadyFollowing = errors.New("already following this user")
)
