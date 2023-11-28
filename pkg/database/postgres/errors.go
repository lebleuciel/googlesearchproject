package postgres

import "github.com/pkg/errors"

var ErrUserWithEmailExist = errors.New("User with the same given email exists")
