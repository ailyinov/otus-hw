package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(in string) (string, error) {
	out := strings.Builder{}
	var prevChar rune

	for _, ch := range in {
		if prevChar == 0 && !unicode.IsLetter(ch) {
			return "", ErrInvalidString
		}
		if unicode.IsDigit(ch) {
			if err := validateNumMod(prevChar); err != nil {
				return "", err
			}
			if err := unpackLetter(prevChar, ch, &out); err != nil {
				return "", err
			}
			prevChar = ch
			continue
		}
		writeChar(prevChar, &out)
		prevChar = ch
	}
	writeChar(prevChar, &out)
	return out.String(), nil
}

func validateNumMod(prevChar rune) error {
	if prevChar == 0 || unicode.IsDigit(prevChar) {
		return ErrInvalidString
	}
	return nil
}

func writeChar(ch rune, sb *strings.Builder) {
	if ch > 0 && !unicode.IsDigit(ch) {
		sb.WriteString(string(ch))
	}
}

func unpackLetter(ch rune, num rune, sb *strings.Builder) error {
	if n, err := strconv.Atoi(string(num)); err == nil {
		sb.WriteString(strings.Repeat(string(ch), n))
		return nil
	}
	return ErrInvalidString
}
