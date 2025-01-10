package helpers

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strings"
	"unicode"
)

func SnakeToCamel(snake string) string {
	parts := strings.Split(snake, "_")
	var result strings.Builder
	for _, part := range parts {
		result.WriteString(cases.Title(language.English).String(part))
	}
	return result.String()
}

func SplitCamelCase(s string) []string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	split := re.ReplaceAllString(s, "${1} ${2}")
	return strings.Split(split, " ")
}

func ToLowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

func CamelToSnake(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
