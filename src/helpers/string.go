package helpers

import (
	"errors"
)

const (
	WHITESPACE_CUTSET = " \t\n"
)

func SubString(text string, start, end int) string {
	textRunes := []rune(text)
	startPos := start
	endPos := end
	if start < 0 {
		startPos = len(textRunes) - (start * -1)
	}
	if end > len(textRunes) {
		endPos = len(textRunes)
	}
	if end < 0 {
		endPos = len(textRunes) - (end * -1)
	}
	if endPos-startPos <= 0 {
		return ""
	}
	return string(textRunes[startPos:endPos])
}

// Retrieves a comma-separated list of values from a string
func GetCSV(text string) ([]string, error) {
	textRunes := []rune(text)
	inQuotes := false
	n := 0
	for i := 0; i < len(textRunes); i++ {
		c := textRunes[i]
		if c == '"' {
			if !inQuotes {
				inQuotes = true
			} else if inQuotes {
				if textRunes[i-1] == '\\' {
					continue
				}
				inQuotes = false
				if i == len(textRunes)-1 {
					n++
				}
			}
		} else if !inQuotes && (c == ',' || i == len(textRunes)-1) {
			n++
		}
	}
	if inQuotes {
		return nil, errors.New("Inconsistent double quotes in text")
	}
	values := make([]string, n)
	n = 0
	a, b := -1, -1
	inQuotes = false
	skip := false
	for i := 0; i < len(textRunes); i++ {
		c := textRunes[i]
		if skip {
			skip = false
			continue
		}
		if c != ' ' {
			if a < 0 {
				a = i
			}
			b = i
		}
		if c == '"' {
			if !inQuotes {
				inQuotes = true
				a = i + 1
			} else if inQuotes {
				if textRunes[i-1] == '\\' {
					continue
				}
				inQuotes = false
				values[n] = string(textRunes[a:b])
				n++
				a, b = -1, -1
				if i == len(textRunes)-1 {
					values[n] = string(textRunes[a:b])
					n++
				}
				skip = true
			}
		} else if !inQuotes && a >= 0 && (c == ',' || i == len(textRunes)-1) {
			if i == len(textRunes)-1 && c != ',' {
				b++
			}
			values[n] = string(textRunes[a:b])
			n++
			a, b = -1, -1
		}
	}
	return values, nil
}
