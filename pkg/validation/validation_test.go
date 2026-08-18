package validation

import "testing"

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "example from iCRM", value: "test@google.com", want: true},
		{name: "standard email", value: "alice@example.com", want: true},
		{name: "plus tag", value: "alice+tag@example.com", want: true},
		{name: "subdomain tld", value: "user@mail.google.com", want: true},
		{name: "missing at", value: "alice.example.com", want: false},
		{name: "missing domain", value: "alice@", want: false},
		{name: "missing local", value: "@example.com", want: false},
		{name: "spaces", value: "alice @example.com", want: false},
		{name: "empty", value: "", want: false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := IsValidEmail(testCase.value)
			if got != testCase.want {
				t.Fatalf("IsValidEmail(%q) = %v, want %v", testCase.value, got, testCase.want)
			}
		})
	}
}

func TestIsValidObjectID(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "lowercase hex", value: "507f1f77bcf86cd799439011", want: true},
		{name: "uppercase hex", value: "507F1F77BCF86CD799439011", want: true},
		{name: "mixed case", value: "507f1F77bcf86CD799439011", want: true},
		{name: "too short", value: "507f1f77bcf86cd79943901", want: false},
		{name: "too long", value: "507f1f77bcf86cd799439011a", want: false},
		{name: "non hex", value: "not-an-object-id-xxxxxx", want: false},
		{name: "empty", value: "", want: false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := IsValidObjectID(testCase.value)
			if got != testCase.want {
				t.Fatalf("IsValidObjectID(%q) = %v, want %v", testCase.value, got, testCase.want)
			}
		})
	}
}
