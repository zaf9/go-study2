package password

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestValidate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cases := []struct {
			name    string
			raw     string
			wantErr bool
		}{
			{name: "valid", raw: "Aa1!aaaa", wantErr: false},
			{name: "missing_upper", raw: "aa1!aaaa", wantErr: true},
			{name: "missing_lower", raw: "AA1!AAAA", wantErr: true},
			{name: "missing_digit", raw: "Aa!aaaaa", wantErr: true},
			{name: "missing_special", raw: "Aa1aaaaa", wantErr: true},
			{name: "too_short", raw: "Aa1!", wantErr: true},
			{name: "empty", raw: "", wantErr: true},
		}

		for _, c := range cases {
			err := Validate(c.raw)
			if c.wantErr {
				t.AssertNE(err, nil)
			} else {
				t.AssertNil(err)
			}
		}
	})
}

func TestHashAndVerify(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		hash, err := Hash("Test1234!")
		t.AssertNil(err)
		t.AssertNE(hash, "")

		t.AssertNil(Verify(hash, "Test1234!"))
		t.AssertNE(Verify(hash, "BadPassword"), nil)
	})
}

func TestHashRejectsInvalidPassword(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		hash, err := Hash("weakpass")
		t.AssertNE(err, nil)
		t.Assert(hash, "")
	})
}
