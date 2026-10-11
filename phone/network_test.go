package phone

import (
	"sync"
	"testing"
)

func TestIdentifyNetwork(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "MTN local number",
			input:    "08031234567",
			expected: "MTN",
		},
		{
			name:     "Glo local number",
			input:    "08051234567",
			expected: "Glo",
		},
		{
			name:     "Airtel local number",
			input:    "08021234567",
			expected: "Airtel",
		},
		{
			name:     "9mobile international number",
			input:    "+2348091234567",
			expected: "9mobile",
		},
		{
			name:     "M-Tel full allocation",
			input:    "08041234567",
			expected: "M-Tel",
		},
		{
			name:     "Mafab partial allocation in range",
			input:    "08010123456",
			expected: "Mafab",
		},
		{
			name:     "Mafab prefix outside allocated subrange",
			input:    "08011234567",
			expected: "Unknown",
		},
		{
			name:     "unknown mobile prefix with valid format",
			input:    "08001234567",
			expected: "Unknown",
		},
		{
			name:     "special service prefix 0700 is unsupported",
			input:    "07001234567",
			expected: "Unknown",
		},
		{
			name:     "withdrawn prefix 0709 is unsupported",
			input:    "07091234567",
			expected: "Unknown",
		},
		{
			name:     "partial allocation prefix 0702 is unsupported",
			input:    "07021234567",
			expected: "Unknown",
		},
		{
			name:     "invalid number",
			input:    "12345",
			expected: "Unknown",
		},
		{
			name:     "vanity reservation 0900 is unsupported",
			input:    "09001234567",
			expected: "Unknown",
		},
		{
			name:     "special service prefix 0800 is unsupported",
			input:    "08001234567",
			expected: "Unknown",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := IdentifyNetwork(test.input)
			if got != test.expected {
				t.Errorf(
					"IdentifyNetwork(%q) = %q; want %q",
					test.input,
					got,
					test.expected,
				)
			}
		})
	}
}

func TestIdentifyNetworkSupportedPrefixes(t *testing.T) {
	// One sample national number per fully allocated 4-digit prefix.
	tests := []struct {
		prefix   string
		expected string
	}{
		{"0703", "MTN"},
		{"0704", "MTN"},
		{"0706", "MTN"},
		{"0707", "MTN"},
		{"0803", "MTN"},
		{"0806", "MTN"},
		{"0810", "MTN"},
		{"0813", "MTN"},
		{"0814", "MTN"},
		{"0816", "MTN"},
		{"0903", "MTN"},
		{"0906", "MTN"},
		{"0913", "MTN"},
		{"0916", "MTN"},

		{"0701", "Airtel"},
		{"0708", "Airtel"},
		{"0802", "Airtel"},
		{"0808", "Airtel"},
		{"0812", "Airtel"},
		{"0901", "Airtel"},
		{"0902", "Airtel"},
		{"0904", "Airtel"},
		{"0907", "Airtel"},
		{"0911", "Airtel"},
		{"0912", "Airtel"},

		{"0705", "Glo"},
		{"0805", "Glo"},
		{"0807", "Glo"},
		{"0811", "Glo"},
		{"0815", "Glo"},
		{"0905", "Glo"},
		{"0915", "Glo"},

		{"0809", "9mobile"},
		{"0817", "9mobile"},
		{"0818", "9mobile"},
		{"0908", "9mobile"},
		{"0909", "9mobile"},

		{"0804", "M-Tel"},
	}

	for _, test := range tests {
		t.Run(test.prefix, func(t *testing.T) {
			local := test.prefix + "1234567"
			international := "+234" + test.prefix[1:] + "1234567"

			if got := IdentifyNetwork(local); got != test.expected {
				t.Errorf("IdentifyNetwork(%q) = %q; want %q", local, got, test.expected)
			}
			if got := IdentifyNetwork(international); got != test.expected {
				t.Errorf("IdentifyNetwork(%q) = %q; want %q", international, got, test.expected)
			}

			if mapped, ok := networkPrefixes[test.prefix]; !ok || mapped != test.expected {
				t.Errorf("networkPrefixes[%q] = %q, ok=%v; want %q", test.prefix, mapped, ok, test.expected)
			}
		})
	}
}

func TestIdentifyNetworkConcurrent(t *testing.T) {
	inputs := []string{
		"08031234567",
		"+2348051234567",
		"08021234567",
		"08091234567",
		"08001234567",
		"08010123456",
	}

	var waitGroup sync.WaitGroup
	for i := 0; i < 32; i++ {
		waitGroup.Add(1)
		go func(index int) {
			defer waitGroup.Done()
			input := inputs[index%len(inputs)]
			_ = IdentifyNetwork(input)
		}(i)
	}
	waitGroup.Wait()
}
