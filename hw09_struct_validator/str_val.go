package hw09structvalidator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type StrVal struct {
	field     string
	value     string
	ruleValue string
	valid     func() (bool, error)
}

func NewStrVal(field string, value string, rule *ValidationRule) *StrVal {
	strVal := &StrVal{
		field:     field,
		value:     value,
		ruleValue: rule.Value,
	}
	switch rule.Name {
	case "len":
		strVal.valid = strVal.len
	case "regexp":
		strVal.valid = strVal.regexp
	case "in":
		strVal.valid = strVal.in
	}
	return strVal
}

func (s *StrVal) Valid() (bool, error) {
	return s.valid()
}

func (s *StrVal) Field() string {
	return s.field
}

func (s *StrVal) len() (bool, error) {
	if v, err := strconv.Atoi(s.ruleValue); err != nil {
		return false, err
	} else {
		if ok := len(s.value) == v; !ok {
			return ok, fmt.Errorf("value length does not match: %w", ErrNotValidValue)
		} else {
			return ok, nil
		}
	}
}

func (s *StrVal) regexp() (bool, error) {
	reg, err := regexp.Compile(s.ruleValue)
	if err != nil {
		return false, err
	}
	if ok := reg.MatchString(s.value); !ok {
		return ok, fmt.Errorf("value does not match the pattern: %w", ErrNotValidValue)
	} else {
		return ok, nil
	}
}

func (s *StrVal) in() (bool, error) {
	values := strings.Split(s.ruleValue, ",")
	for _, v := range values {
		if s.value == v {
			return true, nil
		}
	}
	return false, fmt.Errorf("value out of required set: %w", ErrNotValidValue)
}
