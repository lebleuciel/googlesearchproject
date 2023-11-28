package auth

import "github.com/pkg/errors"

var ErrEmptySecretKey = errors.New("Secret key for auth module should not be empty")
var ErrNilUserRepo = errors.New("User repository should not be nil for auth module creation")
var ErrUserObjectNotFound = errors.New("User object not found in gin context")
var ErrInvalidUserObjectType = errors.New("Could not cast object to user")
