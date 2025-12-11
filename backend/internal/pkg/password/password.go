package password

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordEmpty      = errors.New("密码不能为空")
	ErrPasswordTooShort   = errors.New("密码长度需至少 8 位")
	ErrPasswordComplexity = errors.New("密码需包含大写、小写、数字和特殊字符")
	lowerPattern          = regexp.MustCompile(`[a-z]`)
	upperPattern          = regexp.MustCompile(`[A-Z]`)
	digitPattern          = regexp.MustCompile(`[0-9]`)
	specialCharPattern    = regexp.MustCompile(`[^A-Za-z0-9]`)
)

func Validate(raw string) error {
	if raw == "" {
		return ErrPasswordEmpty
	}
	if len([]rune(raw)) < 8 {
		return ErrPasswordTooShort
	}
	if !lowerPattern.MatchString(raw) ||
		!upperPattern.MatchString(raw) ||
		!digitPattern.MatchString(raw) ||
		!specialCharPattern.MatchString(raw) {
		return ErrPasswordComplexity
	}
	return nil
}

// Hash 对密码进行 bcrypt 哈希。
func Hash(raw string) (string, error) {
	if err := Validate(raw); err != nil {
		return "", err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// Verify 校验密码是否匹配哈希。
func Verify(hashed string, raw string) error {
	if hashed == "" {
		return errors.New("哈希值不能为空")
	}
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
}
