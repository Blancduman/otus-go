package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

var reg = regexp.MustCompile(`(\\.|.)([0-9]*)?`)

func Unpack(str string) (string, error) {
	var result strings.Builder

	if str == "" {
		return "", nil
	}

	for _, pair := range reg.FindAllStringSubmatch(str, -1) {
		parsedAmount := 1
		char, amount := pair[1], pair[2]

		if len(amount) != 0 {
			if v, err := strconv.Atoi(amount); err == nil {
				parsedAmount = v
			}
		}

		if len(char) > 1 && strings.Contains(char, "\\") {
			if unicode.IsDigit(rune(char[1])) || char[1] == '\\' {
				char = string(rune(char[1]))
			}
		}

		if parsedAmount > 9 {
			return "", ErrInvalidString
		}

		result.WriteString(strings.Repeat(char, parsedAmount))
	}

	return result.String(), nil
}
