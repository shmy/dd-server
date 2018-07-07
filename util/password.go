package util

import (
	"golang.org/x/crypto/scrypt"
	"unsafe"
)

const (
	SALT = "@#$%^&*()"
)

// 生成专家级密码 http://wiki.jikexueyuan.com/project/go-web-programming/09.5.html
func GenerateThePassword(password string) *string {
	dk, err := scrypt.Key([]byte(password), []byte(SALT), 16384, 8, 1, 32)
	if err != nil {
		return nil
	}
	return (*string)(unsafe.Pointer(&dk))
}
