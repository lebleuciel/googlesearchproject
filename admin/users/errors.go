package users

import "github.com/pkg/errors"

var ErrNilUserRepo = errors.New("User repository should not be nil")
var ErrNilAuthModule = errors.New("Auth module should not be nil")
var ErrNilUserService = errors.New("User service should not be nil")
