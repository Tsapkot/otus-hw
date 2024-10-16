package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	// Place your code here.
	var builder strings.Builder
	var err error
	var runes = []rune(input)

	var previous rune
	for i := 0; i < len(runes); i++ {
		var current = runes[i]
		var intVal, _ = strconv.Atoi(string(current))

		if unicode.IsDigit(current) {
			if i == 0 || unicode.IsDigit(previous) {
				err = ErrInvalidString
			} else {
				builder.WriteString(strings.Repeat(string(previous), intVal))
			}
		} else if !unicode.IsDigit(current) && i != 0 {
			if !unicode.IsDigit(previous) {
				builder.WriteString(string(previous))
			}
			if i == len(runes)-1 {
				builder.WriteString(string(current))
			}
		}
		previous = current
	}
	return builder.String(), err
}
