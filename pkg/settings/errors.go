package settings

import "github.com/pkg/errors"

var ErrSettingNameEmpty = errors.New("global.name field is required.")
var ErrSettingDuplicatedServerPorts = errors.New("duplicated ports has been found: port number fields in setting.yml should have different values.")
var ErrSettingInvalidEnvironment = errors.New("configs.environment field value is invalid.")
