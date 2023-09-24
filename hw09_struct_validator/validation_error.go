package hw09structvalidator

import (
	"errors"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	strBuilder := strings.Builder{}

	for _, ve := range v {
		errMessage := ve.Err.Error()
		if strBuilder.Len() != 0 {
			errMessage += "\n\t"
		}

		strBuilder.WriteString(errMessage)
	}

	return strBuilder.String()
}

var (
	ErrValNotStruct                = errors.New("provided value is not a struct")
	ErrValIncorrectRule            = errors.New("unable to parse rule")
	ErrValIncorrectStringLength    = errors.New("incorrect string length")
	ErrValRegexString              = errors.New("string do not satisfy regex")
	ErrValNoMatchingElementInSlice = errors.New("there is no matching element")
	ErrValIntLessThanMin           = errors.New("provided int value is less than min")
	ErrValIntMoreThanMax           = errors.New("provided int value is more than max")
	ErrValUnsupportedType          = errors.New("provided type is not supported")
	ErrValTagValueShouldBeInteger  = errors.New("validation tag value should be integer")
	ErrValIncorrectTagRegexPattern = errors.New("provided regexp pattern is invalid")
)

func (v ValidationErrors) Is(target error) bool {
	var tErr ValidationErrors

	if !errors.As(target, &tErr) {
		return false
	}

	if len(v) != len(tErr) {
		return false
	}

	for i := 0; i < len(tErr); i++ {
		areFieldsEqual := v[i].Field == tErr[i].Field
		areErrsEqual := errors.Is(v[i].Err, tErr[i].Err)
		if !areFieldsEqual || !areErrsEqual {
			return false
		}
	}

	return true
}
