package hw09structvalidator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Slice struct {
	validators []Validator
	field      string
}

func NewSlice(field string, values []any, rule *ValidationRule) Validator {
	validators := make([]Validator, 0, len(values))
	switch v := values[0].(type) {
	case int:
		for i, val := range values {
			validators = append(validators, NewIntVal(strconv.Itoa(i), val.(int), rule))
		}
	case string:
		for i, val := range values {
			validators = append(validators, NewStrVal(strconv.Itoa(i), val.(string), rule))
		}
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}

	return &Slice{validators: validators, field: field}
}

func (i *Slice) Valid() (bool, error) {
	okAll := true
	var errs []string
	for _, v := range i.validators {
		_, err := v.Valid()
		if err != nil && errors.Is(err, ErrNotValidValue) {
			okAll = false
			errs = append(errs, fmt.Sprintf("[%s, %v]", v.Field(), err))
		} else if err != nil {
			return false, err
		}
	}
	if !okAll {
		return okAll, fmt.Errorf("int slise contains errors: %s", strings.Join(errs, ", "))
	}
	return okAll, nil
}

func (i *Slice) Field() string {
	return i.field
}
