package mapper

import "unicode"

// ToSnakeCase converts PascalCase/camelCase to snake_case.
func ToSnakeCase(value string) string {
	if value == "" {
		return value
	}
	runes := []rune(value)
	out := make([]rune, 0, len(runes)+4)
	for i, r := range runes {
		isUpper := unicode.IsUpper(r)
		if isUpper {
			if i > 0 {
				prev := runes[i-1]
				nextIsLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
				if unicode.IsLower(prev) || unicode.IsDigit(prev) || nextIsLower {
					out = append(out, '_')
				}
			}
			out = append(out, unicode.ToLower(r))
			continue
		}
		out = append(out, r)
	}
	return string(out)
}
