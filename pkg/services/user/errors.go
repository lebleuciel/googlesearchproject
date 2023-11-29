package user

import "github.com/pkg/errors"

var ErrNilUserRepo = errors.New("User repository can not be nil")
