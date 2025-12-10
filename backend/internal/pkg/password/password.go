package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Hash 对密码进行 bcrypt 哈希。
func Hash(raw string) (string, error) {
	if raw == "" {
		return "", errors.New("密码不能为空")
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
