package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	tagValidate = "validate"
)

var (
	ErrNoValidator   = errors.New("no suitable validator found for field")
	ErrNotValidValue = errors.New("value not valid")
)

type VarType interface {
	int | string
}

type Validator interface {
	Valid() (bool, error)
	Field() string
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationRule struct {
	Name  string
	Value string
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	sb := strings.Builder{}
	for _, e := range v {
		sb.WriteString(fmt.Sprintf("%s: %v \n", e.Field, e.Err))
	}
	return sb.String()
}

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("value must be a pointer")
	}
	value = value.Elem().Elem()
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("value must be a struct")
	}

	valueType := value.Type()

	validators := make([]Validator, 0)
	var valErrs ValidationErrors
	for i := 0; i < value.NumField(); i++ {
		fieldTag, ok := valueType.Field(i).Tag.Lookup(tagValidate)
		if !ok {
			continue
		}

		fieldName := valueType.Field(i).Name
		fieldValue := value.FieldByName(fieldName)
		rules := strings.Split(fieldTag, "|")
		switch fieldValue.Kind() {
		case reflect.Int:
			for _, r := range rules {
				validators = append(validators, NewIntVal(fieldName, int(fieldValue.Int()), NewValidationRule(r)))
			}
		case reflect.String:
			for _, r := range rules {
				validators = append(validators, NewStrVal(fieldName, fieldValue.String(), NewValidationRule(r)))
			}
		case reflect.Array | reflect.Slice:
			sliceKind := valueType.Field(i).Type.Elem().Kind()
			//for _, r := range rules {
			//	validators = append(validators, NewSlice(fieldName, fieldValue, NewValidationRule(r)))
			//}
			if sliceKind == reflect.String {
				for _, r := range rules {
					validators = append(validators, NewStrSlice(fieldName, fieldValue.Interface().([]string), NewValidationRule(r)))
				}
			}
			if sliceKind == reflect.Int {
				for _, r := range rules {
					validators = append(validators, NewIntSlice(fieldName, fieldValue.Interface().([]int), NewValidationRule(r)))
				}
			}
		default:
			valErrs = append(valErrs, ValidationError{
				Field: fieldName,
				Err:   ErrNoValidator,
			})
		}

	}

	for _, v := range validators {
		valid, err := v.Valid()
		if err != nil && !errors.Is(err, ErrNotValidValue) {
			return err
		}
		if !valid {
			valErrs = append(valErrs, ValidationError{
				Field: v.Field(),
				Err:   err,
			})
		}
	}
	return valErrs
}

func NewValidationRule(rule string) *ValidationRule {
	r := strings.Split(rule, ":")
	rName := r[0]
	rVal := r[1]
	return &ValidationRule{
		Name:  rName,
		Value: rVal,
	}
}

//func (r ValidationRule) IntVal(value int) *validator.IntVal {
//	intVal := &validator.IntVal{
//		value:     value,
//		ruleValue: r.value,
//	}
//	switch r.Name {
//	case "min":
//		intVal.valid = intVal.min
//	case "max":
//		intVal.valid = intVal.max
//	case "in":
//		intVal.valid = intVal.in
//
//	}
//	return intVal
//}
