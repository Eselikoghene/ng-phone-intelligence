package phone

import "testing"

func BenchmarkNormalizePhoneNumberLocal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NormalizePhoneNumber("08031234567")
	}
}

func BenchmarkNormalizePhoneNumberInternational(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NormalizePhoneNumber("+2348031234567")
	}
}

func BenchmarkIdentifyNetwork(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IdentifyNetwork("08031234567")
	}
}
