package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(pwdStr string) (encryptPwdStr string, err error) {
	pwd := []byte(pwdStr)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return
	}
	encryptPwdStr = string(hash)
	return
}

func VerifyPassword(encryptPwdStr string, pwdStr string) bool {
	byteEncrypt := []byte(encryptPwdStr)
	bytePwd := []byte(pwdStr)
	err := bcrypt.CompareHashAndPassword(byteEncrypt, bytePwd)
	if err != nil {
		return false
	}
	return true
}
