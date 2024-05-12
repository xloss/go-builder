package builder

import "testing"

func TestRandStr(t *testing.T) {
	var (
		str1 = randStr()
		str2 = randStr()
		str3 = randStr()
	)

	if str1 == str2 {
		t.Error("str1 and str2 should not be equal")
	}
	if str1 == str3 {
		t.Error("str1 and str3 should not be equal")
	}
	if str2 == str3 {
		t.Error("str2 and str3 should not be equal")
	}

	if len(str1) != len(str2) {
		t.Error("len(str1) != len(str2)")
	}
	if len(str1) != len(str3) {
		t.Error("len(str1) != len(str3)")
	}
	if len(str2) != len(str3) {
		t.Error("len(str2) != len(str3)")
	}

	if len(str1) != randStrLen {
		t.Errorf("len(str1) != %d", randStrLen)
	}
}
