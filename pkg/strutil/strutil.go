package strutil

import (
	"strings"
	"unicode"
)

// ToPascalCase converts a string (snake, kebab, camel) to PascalCase (e.g. "user_profile" -> "UserProfile").
func ToPascalCase(s string) string {
	if s == "" {
		return ""
	}
	words := splitIntoWords(s)
	for i, w := range words {
		words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
	}
	return strings.Join(words, "")
}

// ToCamelCase converts a string to camelCase (e.g. "UserProfile" -> "userProfile").
func ToCamelCase(s string) string {
	pascal := ToPascalCase(s)
	if pascal == "" {
		return ""
	}
	return strings.ToLower(pascal[:1]) + pascal[1:]
}

// ToSnakeCase converts a string to snake_case (e.g. "UserProfile" -> "user_profile").
func ToSnakeCase(s string) string {
	words := splitIntoWords(s)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "_")
}

// ToKebabCase converts a string to kebab-case (e.g. "UserProfile" -> "user-profile").
func ToKebabCase(s string) string {
	words := splitIntoWords(s)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "-")
}

// Pluralize gives a basic English plural form.
func Pluralize(s string) string {
	lower := strings.ToLower(s)
	if strings.HasSuffix(lower, "y") && !strings.HasSuffix(lower, "ay") && !strings.HasSuffix(lower, "ey") && !strings.HasSuffix(lower, "oy") && !strings.HasSuffix(lower, "uy") {
		return s[:len(s)-1] + "ies"
	}
	if strings.HasSuffix(lower, "s") || strings.HasSuffix(lower, "x") || strings.HasSuffix(lower, "z") || strings.HasSuffix(lower, "ch") || strings.HasSuffix(lower, "sh") {
		return s + "es"
	}
	return s + "s"
}

func splitIntoWords(s string) []string {
	var words []string
	var current strings.Builder

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '_' || r == '-' || r == ' ' || r == '.' || r == '/' {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
			continue
		}

		if unicode.IsUpper(r) {
			if current.Len() > 0 {
				// Check if previous was upper and next is lower (e.g. "HTTPRequest" -> "HTTP", "Request")
				prevUpper := i > 0 && unicode.IsUpper(runes[i-1])
				nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
				if !prevUpper || nextLower {
					words = append(words, current.String())
					current.Reset()
				}
			}
		}
		current.WriteRune(r)
	}
	if current.Len() > 0 {
		words = append(words, current.String())
	}
	return words
}
