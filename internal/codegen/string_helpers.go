package codegen

import "strings"

func snakeToCamel(snake string) string {
	parts := strings.Split(snake, "_")
	var result strings.Builder
	for _, part := range parts {
		result.WriteString(strings.Title(part))
	}
	return result.String()
}
