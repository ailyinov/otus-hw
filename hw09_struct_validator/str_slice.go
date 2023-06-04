package hw09structvalidator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type StrSlice struct {
	validators []*StrVal
	field      string
}

func NewStrSlice(field string, values []string, rule *ValidationRule) *StrSlice {
	validators := make([]*StrVal, 0, len(values))
	for i, v := range values {
		validators = append(validators, NewStrVal(strconv.Itoa(i), v, rule))
	}
	return &StrSlice{validators: validators, field: field}
}

func (s *StrSlice) Valid() (bool, error) {
	okAll := true
	var errs []string
	for _, v := range s.validators {
		_, err := v.valid()
		if err != nil && errors.Is(err, ErrNotValidValue) {
			okAll = false
			errs = append(errs, fmt.Sprintf("[%s, %v]", v.field, err))
		} else if err != nil {
			return false, err
		}
	}
	if !okAll {
		return okAll, fmt.Errorf("int slise contains errors: %s", strings.Join(errs, ", "))
	}
	return okAll, nil
}

func (s *StrSlice) Field() string {
	return s.field
}
