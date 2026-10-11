package phone

import (
	"errors"
	"testing"
)

func TestNormalizePhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  error
	}{
		{
			name:     "valid local format",
			input:    "08031234567",
			expected: "+2348031234567",
			wantErr:  nil,
		},
		{
			name:     "valid international format with plus",
			input:    "+2348031234567",
			expected: "+2348031234567",
			wantErr:  nil,
		},
		{
			name:     "leading and trailing whitespace local",
			input:    "  08031234567  ",
			expected: "+2348031234567",
			wantErr:  nil,
		},
		{
			name:     "leading and trailing whitespace international",
			input:    "  +2348031234567  ",
			expected: "+2348031234567",
			wantErr:  nil,
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
			wantErr:  ErrEmptyInput,
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: "",
			wantErr:  ErrEmptyInput,
		},
		{
			name:     "too short",
			input:    "080312345",
			expected: "",
			wantErr:  ErrInvalidFormat,
		},
		{
			name:     "too long local",
			input:    "080312345678",
			expected: "",
			wantErr:  ErrInvalidFormat,
		},
		{
			name:     "too long international",
			input:    "+23480312345678",
			expected: "",
			wantErr:  ErrInvalidFormat,
		},
		{
			name:     "letters in local number",
			input:    "0803123456A",
			expected: "",
			wantErr:  ErrInvalidFormat,
		},
		{
			name:     "letters in international number",
			input:    "+234803123456A",
			expected: "",
			wantErr:  ErrInvalidFormat,
		},
		{
			name:     "unsupported characters",
			input:    "0803-123-4567",
			expected: "",
			wantErr:  ErrInvalidFormat,
		},
		{
			name:     "unsupported country code",
			input:    "+18031234567",
			expected: "",
			wantErr:  ErrUnsupportedCountryCode,
		},
		{
			name:     "incorrect country code and national number combination keeps local zero",
			input:    "+23408031234567",
			expected: "",
			wantErr:  ErrInvalidFormat,
		},
		{
			name:     "bare international without plus is not accepted",
			input:    "2348031234567",
			expected: "",
			wantErr:  ErrInvalidFormat,
		},
		{
			name:     "valid format with unrecognized mobile prefix still normalizes",
			input:    "08001234567",
			expected: "+2348001234567",
			wantErr:  nil,
		},
		{
			name:     "unicode digits are rejected",
			input:    "０８０３１２３４５６７",
			expected: "",
			wantErr:  ErrInvalidFormat,
		},
		{
			name:     "internal whitespace is rejected",
			input:    "0803 123 4567",
			expected: "",
			wantErr:  ErrInvalidFormat,
		},
		{
			name:     "plus with incomplete country code",
			input:    "+23",
			expected: "",
			wantErr:  ErrUnsupportedCountryCode,
		},
		{
			name:     "plus 234 with truncated national number",
			input:    "+234803123456",
			expected: "",
			wantErr:  ErrInvalidFormat,
		},
		{
			name:     "tab and newline whitespace around valid local number",
			input:    "\t08031234567\n",
			expected: "+2348031234567",
			wantErr:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NormalizePhoneNumber(test.input)

			if got != test.expected {
				t.Errorf("NormalizePhoneNumber(%q) = %q; want %q", test.input, got, test.expected)
			}

			if test.wantErr == nil {
				if err != nil {
					t.Errorf("NormalizePhoneNumber(%q) unexpected error: %v", test.input, err)
				}
				return
			}

			if !errors.Is(err, test.wantErr) {
				t.Errorf("NormalizePhoneNumber(%q) error = %v; want %v", test.input, err, test.wantErr)
			}
		})
	}
}
