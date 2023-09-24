package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	validationTagName = "validate"
)

type (
	ValidatorMap map[string]string
	ValidatorFn  func(reflect.Value, ValidatorMap) error
)

func Validate(v interface{}) error {
	var errs ValidationErrors

	reflType := reflect.TypeOf(v)
	reflValue := reflect.ValueOf(v)

	if reflType.Kind() != reflect.Struct {
		return ErrValNotStruct
	}

	for i := 0; i < reflValue.NumField(); i++ {
		field := reflType.Field(i)
		value := reflValue.Field(i)

		tag, ok := field.Tag.Lookup(validationTagName)
		if !ok {
			continue
		}

		validator, err := parseTag(tag)
		if err != nil {
			return append(errs, ValidationError{
				Field: field.Name,
				Err:   err,
			})
		}

		errs = append(errs, validate(field, value, validator)...)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validate(field reflect.StructField, value reflect.Value, rules ValidatorMap) ValidationErrors {
	var errs ValidationErrors

	kind := field.Type.Kind()

	switch kind { //nolint: exhaustive
	case reflect.Int:
		if err := validateInt(value, rules); err != nil {
			errs = append(errs, ValidationError{
				Field: field.Name,
				Err:   err,
			})
		}
	case reflect.String:
		if err := validateString(value, rules); err != nil {
			errs = append(errs, ValidationError{
				Field: field.Name,
				Err:   err,
			})
		}
	case reflect.Slice:
		elKind := field.Type.Elem().Kind()
		switch elKind { //nolint: exhaustive
		case reflect.Int:
			for _, err := range validateSlice(validateInt, value, rules) {
				errs = append(errs, ValidationError{
					Field: field.Name,
					Err:   err,
				})
			}
		case reflect.String:
			for _, err := range validateSlice(validateString, value, rules) {
				errs = append(errs, ValidationError{
					Field: field.Name,
					Err:   err,
				})
			}
		default:
			errs = append(errs, ValidationError{
				Field: field.Name,
				Err:   ErrValUnsupportedType,
			})
		}
	default:
		errs = append(errs, ValidationError{
			Field: field.Name,
			Err:   ErrValUnsupportedType,
		})
	}

	return errs
}

func parseTag(tag string) (ValidatorMap, error) {
	vm := make(ValidatorMap)

	for _, rule := range strings.Split(tag, "|") {
		fieldRule := strings.Split(rule, ":")

		if len(fieldRule) != 2 || (fieldRule[0] == "") || fieldRule[1] == "" {
			return nil, fmt.Errorf("%w: rule: %s", ErrValIncorrectRule, tag)
		}

		vm[fieldRule[0]] = fieldRule[1]
	}

	return vm, nil
}

func validateInt(value reflect.Value, rules ValidatorMap) error {
	for k, v := range rules {
		switch k {
		case "min":
			expected, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: invalid value min: %s", ErrValTagValueShouldBeInteger, v)
			}

			if value.Int() < expected {
				return ErrValIntLessThanMin
			}
		case "max":
			expected, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: invalid value max: %s", ErrValTagValueShouldBeInteger, v)
			}

			if value.Int() > expected {
				return ErrValIntMoreThanMax
			}
		case "in":
			for _, integer := range strings.Split(v, ",") {
				expected, err := strconv.ParseInt(integer, 10, 64)
				if err != nil {
					return fmt.Errorf("%w: invalid slice value: %s", ErrValTagValueShouldBeInteger, integer)
				}

				if value.Int() == expected {
					return nil
				}
			}

			return ErrValNoMatchingElementInSlice
		}
	}

	return nil
}

func validateString(value reflect.Value, rules ValidatorMap) error {
	for k, v := range rules {
		switch k {
		case "len":
			expected, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("%w: invalid value length: %s", err, v)
			}

			if value.Len() != expected {
				return ErrValIncorrectStringLength
			}
		case "regexp":
			rx, err := regexp.Compile(v)
			if err != nil {
				return fmt.Errorf("%w; invalid regex: %s", ErrValIncorrectTagRegexPattern, v)
			}

			if !rx.MatchString(value.String()) {
				return ErrValRegexString
			}
		case "in":
			for _, word := range strings.Split(v, ",") {
				if value.String() == word {
					return nil
				}
			}

			return ErrValNoMatchingElementInSlice
		}
	}

	return nil
}

func validateSlice(fn ValidatorFn, values reflect.Value, rules ValidatorMap) []error {
	var errs []error

	for i := 0; i < values.Len(); i++ {
		value := values.Index(i)

		if err := fn(value, rules); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
