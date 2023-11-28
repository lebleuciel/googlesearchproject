package user

import "github.com/pkg/errors"

var ErrNilUserDatabase = errors.New("User database should not be nil")
var ErrGetUserByEmail = errors.New("Could not get user with given email address")
