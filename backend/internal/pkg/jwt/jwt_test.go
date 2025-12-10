package jwt

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestGenerateAndVerify(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := Configure(Options{
			Secret:             "0123456789abcdef",
			Issuer:             "test-issuer",
			AccessTokenExpiry:  time.Hour,
			RefreshTokenExpiry: 2 * time.Hour,
		})
		t.AssertNil(err)

		token, err := GenerateAccessToken(1001)
		t.AssertNil(err)
		t.AssertNE(token, "")

		claims, err := VerifyToken(token)
		t.AssertNil(err)
		t.Assert(claims.UserID, int64(1001))
		t.Assert(claims.Issuer, "test-issuer")
	})
}
