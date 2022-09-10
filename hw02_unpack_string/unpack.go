package hw02unpackstring

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(in string) (string, error) {
	r := strings.NewReader(in)
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanRunes)
	sb := strings.Builder{}
	var prevChar string
	for s.Scan() {
		ch := s.Text()
		if len(prevChar) == 0 && !isLetter(ch) {
			return "", ErrInvalidString
		}
		if isNumeric(ch) {
			if err := validateNumMod(prevChar); err != nil {
				return "", err
			}
			if err := unpackLetter(prevChar, ch, &sb); err != nil {
				return "", err
			}
			prevChar = ch
			continue
		}
		writeChar(prevChar, &sb)
		prevChar = ch
	}
	writeChar(prevChar, &sb)
	return sb.String(), nil
}

func validateNumMod(prevChar string) error {
	if len(prevChar) == 0 || isNumeric(prevChar) {
		return ErrInvalidString
	}
	return nil
}

func writeChar(prevChar string, sb *strings.Builder) {
	if len(prevChar) > 0 && !isNumeric(prevChar) {
		sb.WriteString(prevChar)
	}
}

func unpackLetter(ch string, num string, sb *strings.Builder) error {
	if n, err := strconv.Atoi(num); err == nil {
		sb.WriteString(strings.Repeat(ch, n))
		return nil
	}
	return ErrInvalidString
}

func isNumeric(ch string) bool {
	runes := []rune(ch)
	return unicode.IsDigit(runes[0])
}

func isLetter(ch string) bool {
	runes := []rune(ch)
	return unicode.IsLetter(runes[0])
}
