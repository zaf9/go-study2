package password

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestHashAndVerify(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		hash, err := Hash("Test1234!")
		t.AssertNil(err)
		t.AssertNE(hash, "")

		t.AssertNil(Verify(hash, "Test1234!"))
		t.AssertNE(Verify(hash, "BadPassword"), nil)
	})
}
