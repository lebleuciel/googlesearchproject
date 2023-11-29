package forwarder

import "github.com/pkg/errors"

var ErrNilAuthModule = errors.New("Auth module should not be nil")
var ErrEmptyStoreHost = errors.New("Store host should not be empty")
var ErrEmptyAdminPort = errors.New("Admin Port should not be empty")
var ErrEmptyBackendPort = errors.New("Backend Port should not be empty")
var ErrEmptyUserHeaderKey = errors.New("UserHeaderKey should not be empty")
