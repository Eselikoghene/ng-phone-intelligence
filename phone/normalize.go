package phone

import (
	"errors"
	"strings"
)

// Input policy for NormalizePhoneNumber:
//   - Leading and trailing whitespace is trimmed.
//   - Accepted formats:
//     1. Nigerian local mobile form: 0XXXXXXXXXX (11 digits, leading 0).
//     2. E.164 international form with plus: +234XXXXXXXXXX (14 characters).
//   - Bare international form without the leading plus (for example 2348031234567)
//     is not accepted.
//   - Digits only after trimming; letters and other characters are rejected.
//   - Only the Nigerian country code +234 is supported.
//
// A successful result means the input has a valid Nigerian mobile number shape.
// It does not mean the number is assigned, active, or currently served by any network.

var (
	// ErrEmptyInput is returned when the input is empty or only whitespace.
	ErrEmptyInput = errors.New("phone number is empty")

	// ErrInvalidFormat is returned when the input is not a supported Nigerian
	// local or +234 international number format.
	ErrInvalidFormat = errors.New("phone number format is invalid")

	// ErrUnsupportedCountryCode is returned when the input uses a non-Nigerian
	// E.164 country code.
	ErrUnsupportedCountryCode = errors.New("unsupported country code")
)

// NormalizePhoneNumber converts a Nigerian phone number to E.164 (+234...) form.
// See the package input policy above for accepted formats.
func NormalizePhoneNumber(phone string) (string, error) {
	phone = strings.TrimSpace(phone)

	if phone == "" {
		return "", ErrEmptyInput
	}

	if strings.HasPrefix(phone, "+") && !strings.HasPrefix(phone, "+234") {
		return "", ErrUnsupportedCountryCode
	}

	if strings.HasPrefix(phone, "0") && len(phone) == 11 {
		if !isDigits(phone) {
			return "", ErrInvalidFormat
		}
		return "+234" + phone[1:], nil
	}

	if strings.HasPrefix(phone, "+234") && len(phone) == 14 {
		if !isDigits(phone[1:]) {
			return "", ErrInvalidFormat
		}
		return phone, nil
	}

	return "", ErrInvalidFormat
}

// isDigits reports whether every byte in value is an ASCII digit ('0'-'9').
// Unicode digits and other numeric runes are rejected.
func isDigits(value string) bool {
	for i := 0; i < len(value); i++ {
		char := value[i]
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
