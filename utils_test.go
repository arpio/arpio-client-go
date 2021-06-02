package arpio

import "testing"

func TestSliceConstainsString(t *testing.T) {
	if !SliceContainsString("b", []string{"a", "b", "c"}) {
		t.Fatalf("should contain b")
	}
	if SliceContainsString("x", []string{"a", "b", "c"}) {
		t.Fatalf("should not contain x")
	}
}

func TestSliceConstainsInt(t *testing.T) {
	if !SliceContainsInt(2, []int{1, 2, 3}) {
		t.Fatalf("should contain 2")
	}
	if SliceContainsInt(9, []int{1, 2, 3}) {
		t.Fatalf("should not contain 9")
	}
}
