package hw09structvalidator

import (
	"fmt"
	"strconv"
	"strings"
)

type IntVal struct {
	field     string
	value     int
	ruleValue string
	valid     func() (bool, error)
}

func NewIntVal(field string, value int, rule *ValidationRule) *IntVal {
	intVal := &IntVal{
		field:     field,
		value:     value,
		ruleValue: rule.Value,
	}
	switch rule.Name {
	case "min":
		intVal.valid = intVal.min
	case "max":
		intVal.valid = intVal.max
	case "in":
		intVal.valid = intVal.in
	default:
		intVal.valid = intVal.err
	}
	return intVal
}

func (i *IntVal) Valid() (bool, error) {
	return i.valid()
}

func (i *IntVal) Field() string {
	return i.field
}

func (i *IntVal) err() (bool, error) {
	return false, ErrNoValidator
}

func (i *IntVal) max() (bool, error) {
	if v, err := strconv.Atoi(i.ruleValue); err != nil {
		return false, err
	} else {
		if ok := i.value <= v; !ok {
			return ok, fmt.Errorf("value exceeded max: %w", ErrNotValidValue)
		} else {
			return ok, nil
		}
	}
}

func (i *IntVal) min() (bool, error) {
	if v, err := strconv.Atoi(i.ruleValue); err != nil {
		return false, err
	} else {
		if ok := i.value >= v; !ok {
			return ok, fmt.Errorf("value less than min: %w", ErrNotValidValue)
		} else {
			return ok, nil
		}
	}
}

func (i *IntVal) in() (bool, error) {
	values := strings.Split(i.ruleValue, ",")
	for _, vStr := range values {
		if v, err := strconv.Atoi(vStr); err != nil {
			return false, err
		} else if i.value == v {
			return true, nil
		}
	}
	return false, fmt.Errorf("value out of required set: %w", ErrNotValidValue)
}
