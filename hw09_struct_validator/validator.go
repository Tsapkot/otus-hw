package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrInput             = errors.New("input must be a struct")
	ErrInvalidRuleFormat = func(rule string) error {
		return fmt.Errorf("invalid rule format: %s", rule)
	}
	ErrUnsupportedType = func(field reflect.Value) error {
		return fmt.Errorf("unsupported type: %v", field.Kind())
	}
	ErrUnsupportedValidator = func(validatorType string) error {
		return fmt.Errorf("unsupported validator: %s", validatorType)
	}
	ErrValueOutOfRange = func(def string, val int) error {
		return fmt.Errorf("value must be %s than or equal to %d", def, val)
	}
	ErrValueOutOfValues = func(set string) error {
		return fmt.Errorf("value must be one of %s", set)
	}
	ErrStringLengthMismatch = func(expectedLen int) error {
		return fmt.Errorf("string length must be exactly %d", expectedLen)
	}
	ErrStringDoesNotMatchPattern = func(pattern string) error {
		return fmt.Errorf("string does not match pattern: %s", pattern)
	}
)

func (v ValidationErrors) Error() string {
	errs := make([]string, len(v))
	for i, err := range v {
		errs[i] = fmt.Sprintf("%s: %s", err.Field, err.Err)
	}

	return strings.Join(errs, ", ")
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	t := reflect.TypeOf(v)

	if val.Kind() != reflect.Struct {
		return ErrInput
	}

	var validationErrors ValidationErrors

	for i := 0; i < t.NumField(); i++ {
		field := val.Field(i)
		structField := t.Field(i)
		tag := structField.Tag.Get("validate")

		if structField.PkgPath != "" || tag == "" {
			continue
		}

		rules := strings.Split(tag, "|")

		for _, rule := range rules {
			err := defineValidation(rule, field)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: structField.Name, Err: err})
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func defineValidation(rule string, field reflect.Value) error {
	parts := strings.SplitN(rule, ":", 2)
	if len(parts) != 2 {
		return ErrInvalidRuleFormat(rule)
	}

	validatorType := parts[0]
	validatorValue := parts[1]

	switch field.Kind() {
	case reflect.Int:
		if err := validateInt(int(field.Int()), validatorType, validatorValue); err != nil {
			return err
		}
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.Int {
			for j := 0; j < field.Len(); j++ {
				if err := validateInt(int(field.Index(j).Int()), validatorType, validatorValue); err != nil {
					return err
				}
			}
		} else if field.Type().Elem().Kind() == reflect.String {
			for j := 0; j < field.Len(); j++ {
				if err := validateString(field.Index(j).String(), validatorType, validatorValue); err != nil {
					return err
				}
			}
		}
	case reflect.String:
		if err := validateString(field.String(), validatorType, validatorValue); err != nil {
			return err
		}
	default:
		return ErrUnsupportedType(field)
	}
	return nil
}

func validateString(value string, validatorType, validatorValue string) error {
	switch validatorType {
	case "len":
		expectedLen, err := strconv.Atoi(validatorValue)
		if err != nil || len(value) != expectedLen {
			return ErrStringLengthMismatch(expectedLen)
		}
	case "regexp":
		pattern := strings.ReplaceAll(validatorValue, "\\d", `\d`)
		re, err := regexp.Compile(pattern)
		if err != nil || !re.MatchString(value) {
			return ErrStringDoesNotMatchPattern(pattern)
		}
	case "in":
		inSet := strings.Split(validatorValue, ",")
		found := false
		for _, v := range inSet {
			if value == v {
				found = true
				break
			}
		}
		if !found {
			return ErrValueOutOfValues(validatorValue)
		}
	default:
		return ErrUnsupportedValidator(validatorType)
	}
	return nil
}

func validateInt(value int, validatorType, validatorValue string) error {
	switch validatorType {
	case "min":
		minVal, err := strconv.Atoi(validatorValue)
		if err != nil || value < minVal {
			return ErrValueOutOfRange("greater", minVal)
		}
	case "max":
		maxVal, err := strconv.Atoi(validatorValue)
		if err != nil || value > maxVal {
			return ErrValueOutOfRange("less", maxVal)
		}
	case "in":
		inSet := strings.Split(validatorValue, ",")
		found := false
		for _, v := range inSet {
			num, err := strconv.Atoi(v)
			if err == nil && value == num {
				found = true
				break
			}
		}
		if !found {
			return ErrValueOutOfValues(validatorValue)
		}
	default:
		return ErrUnsupportedValidator(validatorType)
	}
	return nil
}
