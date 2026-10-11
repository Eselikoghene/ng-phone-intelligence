package phone

import "strings"

// networkPrefixes maps a 4-digit Nigerian national prefix (including the
// leading 0) to the originally allocated mobile operator.
//
// Source: Nigerian Communications Commission Mobile Number Allocation Table
// https://www.ncc.gov.ng/operators/mobile-number-allocation-table
// Allocation data as at December 2023; page last updated 25 November 2024.
// Verified for this project on 11 October 2026.
//
// Prefix lookup identifies the originally allocated operator only. Because of
// mobile number portability, it cannot reliably establish the current serving
// network. A match also does not mean the number is assigned or active.
//
// Partial allocations are not listed here; see identifyPartialAllocation.
var networkPrefixes = map[string]string{
	// MTN (full 000 0000 to 999 9999)
	"0703": "MTN",
	"0704": "MTN",
	"0706": "MTN",
	"0707": "MTN",
	"0803": "MTN",
	"0806": "MTN",
	"0810": "MTN",
	"0813": "MTN",
	"0814": "MTN",
	"0816": "MTN",
	"0903": "MTN",
	"0906": "MTN",
	"0913": "MTN",
	"0916": "MTN",

	// Airtel (full 000 0000 to 999 9999)
	"0701": "Airtel",
	"0708": "Airtel",
	"0802": "Airtel",
	"0808": "Airtel",
	"0812": "Airtel",
	"0901": "Airtel",
	"0902": "Airtel",
	"0904": "Airtel",
	"0907": "Airtel",
	"0911": "Airtel",
	"0912": "Airtel",

	// Glo (full 000 0000 to 999 9999)
	"0705": "Glo",
	"0805": "Glo",
	"0807": "Glo",
	"0811": "Glo",
	"0815": "Glo",
	"0905": "Glo",
	"0915": "Glo",

	// 9mobile / EMTS (full 000 0000 to 999 9999)
	"0809": "9mobile",
	"0817": "9mobile",
	"0818": "9mobile",
	"0908": "9mobile",
	"0909": "9mobile",

	// M-Tel (full 000 0000 to 999 9999)
	"0804": "M-Tel",
}

// IdentifyNetwork returns the originally allocated mobile network for a
// Nigerian number, or "Unknown" when the number is invalid or the prefix is
// not in the verified allocation set.
//
// This is original-allocation lookup only. It does not report the current
// serving operator, assignment status, or whether the line is active.
func IdentifyNetwork(phone string) string {
	normalized, err := NormalizePhoneNumber(phone)
	if err != nil {
		return "Unknown"
	}
	return identifyNetworkNormalized(normalized)
}

// identifyNetworkNormalized assumes normalized is a valid +234 E.164 value.
// networkPrefixes is package-level and treated as immutable after init, so
// concurrent reads from IdentifyNetwork are safe without additional locking.
func identifyNetworkNormalized(normalized string) string {
	// Remove the +234 country code; localNumber is 10 digits (no leading 0).
	if !strings.HasPrefix(normalized, "+234") || len(normalized) != 14 {
		return "Unknown"
	}
	localNumber := normalized[4:]

	if network, ok := identifyPartialAllocation(localNumber); ok {
		return network
	}

	prefix := "0" + localNumber[:3]
	if network, exists := networkPrefixes[prefix]; exists {
		return network
	}

	return "Unknown"
}

// identifyPartialAllocation handles NCC ranges that are not full 10-million
// blocks under a single 4-digit prefix.
//
// 0801 / MAFAB is allocated only for 000 0000 to 099 9999, so only numbers
// whose national significant number starts with 8010 are mapped.
func identifyPartialAllocation(localNumber string) (string, bool) {
	if strings.HasPrefix(localNumber, "8010") {
		return "Mafab", true
	}
	return "", false
}
