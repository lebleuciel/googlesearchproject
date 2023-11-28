package server

import "github.com/pkg/errors"

var ErrNilAuthModule = errors.New("Gateway auth module can not be nil")
var ErrNilFileModule = errors.New("Gateway file module can not be nil")
