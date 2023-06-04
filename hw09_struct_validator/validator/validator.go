package validator

import (
	"errors"
)

var ErrNoValidator = errors.New("no suitable validator found for type")
