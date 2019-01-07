package testpkg

import "testing"

func TestIntAdd(t *testing.T) {
	if IntAdd(1, 2) == 3 {
		t.Log("IntAdd PASS")
	} else {
		t.Error("IntAdd FAIL")
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IntAdd(1, 2)
	}
}
