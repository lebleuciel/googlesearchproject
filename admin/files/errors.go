package files

import "github.com/pkg/errors"

var ErrNilFileRepo = errors.New("File repository should not be nil")
var ErrNilAuthModule = errors.New("Auth module should not be nil")
var ErrNilFileService = errors.New("File service should not be nil")
