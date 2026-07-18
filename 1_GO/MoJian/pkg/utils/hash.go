package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 对明文密码进行 bcrypt 哈希加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 校验明文密码是否与哈希值匹配
// 返回 true 表示匹配，false 表示不匹配
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
