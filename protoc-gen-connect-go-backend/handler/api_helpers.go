package handler

import (
	"fmt"
	"strings"
)

func toApiTarget(packageName string) string {
	parts := strings.Split(packageName, ".")
	return fmt.Sprintf("api-%s", strings.Join(parts, "-"))
}
