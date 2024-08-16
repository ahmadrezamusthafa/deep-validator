package utils

import "testing"

func TestConvertToSnakeCase(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{"SingleWord", "Example", "example"},
		{"CamelCase", "thisIsAnExample", "this_is_an_example"},
		{"PascalCase", "AnotherExampleString", "another_example_string"},
		{"AllLowercase", "alreadylowercase", "alreadylowercase"},
		{"AllUppercase", "ALLUPPERCASE", "alluppercase"},
		{"MixedCaseWithNumbers", "Example123Test", "example123_test"},
		{"LeadingUppercase", "ExampleStartsWithUppercase", "example_starts_with_uppercase"},
		{"ConsecutiveUppercase", "ConsecutiveUPPERCaseLetters", "consecutive_uppercase_letters"},
		{"ConsecutiveUppercaseAbbreviation", "ID", "id"},
		{"EmptyString", "", ""},
		{"SingleCharacter", "A", "a"},
		{"NonAlphaNumeric", "StringWith123NumbersAndSymbols!", "string_with123_numbers_and_symbols!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToSnakeCase(tt.input)
			if result != tt.output {
				t.Errorf("ConvertToSnakeCase(%s) = %s; want %s", tt.input, result, tt.output)
			}
		})
	}
}
