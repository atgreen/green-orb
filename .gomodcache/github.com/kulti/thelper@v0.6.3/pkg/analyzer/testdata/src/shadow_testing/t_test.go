package testing_test

import testing "shadow_testing"

func tOtherPackage(t *testing.T) {
	_ = 0
	t.Helper()
}

func bOtherPackage(b *testing.B) {
	_ = 0
	b.Helper()
}

func tbOtherPackage(tb *testing.B) {
	_ = 0
	tb.Helper()
}
