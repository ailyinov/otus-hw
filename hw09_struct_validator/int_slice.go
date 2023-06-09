package hw09structvalidator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type IntSlice struct {
	validators []*IntVal
	field      string
}

func NewIntSlice(field string, values []int, rule *ValidationRule) *IntSlice {
	validators := make([]*IntVal, 0, len(values))
	for i, v := range values {
		validators = append(validators, NewIntVal(strconv.Itoa(i), v, rule))
	}
	return &IntSlice{validators: validators, field: field}
}

func (i *IntSlice) Valid() (bool, error) {
	okAll := true
	var errs []string
	for _, v := range i.validators {
		_, err := v.valid()
		if err != nil && errors.Is(err, ErrNotValidValue) {
			okAll = false
			errs = append(errs, fmt.Sprintf("[%s, %v]", v.field, err))
		} else if err != nil {
			return false, err
		}
	}
	if !okAll {
		return okAll, fmt.Errorf("%w, int slise contains errors: %s", ErrNotValidValue, strings.Join(errs, ", "))
	}
	return okAll, nil
}

func (i *IntSlice) Field() string {
	return i.field
}
