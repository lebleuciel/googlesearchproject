package server

import "github.com/pkg/errors"

var ErrNilFileModule = errors.New("Admin file module can not be nil")
var ErrNilUserModule = errors.New("Admin user module can not be nil")
