package strutil

import (
	"testing"
)

func TestCases(t *testing.T) {
	tests := []struct {
		input  string
		pascal string
		camel  string
		snake  string
		kebab  string
		plural string
	}{
		{"user", "User", "user", "user", "user", "Users"},
		{"user_profile", "UserProfile", "userProfile", "user_profile", "user-profile", "UserProfiles"},
		{"order_item", "OrderItem", "orderItem", "order_item", "order-item", "OrderItems"},
		{"category", "Category", "category", "category", "category", "Categories"},
	}

	for _, tt := range tests {
		if got := ToPascalCase(tt.input); got != tt.pascal {
			t.Errorf("ToPascalCase(%q) = %q, want %q", tt.input, got, tt.pascal)
		}
		if got := ToCamelCase(tt.input); got != tt.camel {
			t.Errorf("ToCamelCase(%q) = %q, want %q", tt.input, got, tt.camel)
		}
		if got := ToSnakeCase(tt.input); got != tt.snake {
			t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, got, tt.snake)
		}
		if got := ToKebabCase(tt.input); got != tt.kebab {
			t.Errorf("ToKebabCase(%q) = %q, want %q", tt.input, got, tt.kebab)
		}
		if got := Pluralize(tt.pascal); got != tt.plural {
			t.Errorf("Pluralize(%q) = %q, want %q", tt.pascal, got, tt.plural)
		}
	}
}
