package helpers

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"unicode"
)

func ToGoAlias(packageName string) string {
	// Split the string by the dot
	parts := strings.Split(packageName, ".")
	// Join all parts except the last one
	prefix := strings.Join(parts[:len(parts)-1], "")
	// Capitalize the first letter of the last part (version)
	lastPart := []rune(parts[len(parts)-1])
	if len(lastPart) > 0 {
		lastPart[0] = unicode.ToUpper(lastPart[0])
	}
	// Concatenate the prefix with the modified last part
	return prefix + string(lastPart)
}

func ToGoPackageName(packageName string) string {
	parts := strings.Split(packageName, ".")
	return fmt.Sprintf("api_%s", strings.Join(parts, "_"))
}

// ToGoFieldName converts a snake_case string to CamelCase and handles special cases like "id" to "ID".
func ToGoFieldName(snake string) string {
	// Split the string by underscores
	parts := strings.Split(snake, "_")

	// List of special cases like "id" -> "ID"
	specialCases := map[string]string{
		"id": "ID",
	}

	var result strings.Builder
	for _, part := range parts {
		if val, ok := specialCases[part]; ok {
			result.WriteString(val)
		} else {
			// Capitalize the first letter of each part
			result.WriteString(cases.Title(language.English).String(part))
		}
	}
	return result.String()
}
