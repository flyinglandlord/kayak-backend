package utils

import (
	"math/rand"
	"time"
)

func GenerateInvitationCode(n int) string {
	letters := []byte("0123456789qwertyuioplkjhgfdsazxcvbnmMNBVCXZASDFGHJKLPOIUYTREWQ")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[r.Intn(len(letters))]
	}
	return string(result)
}

func GenerateDigitalCode(n int) string {
	letters := []byte("0123456789")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[r.Intn(len(letters))]
	}
	return string(result)
}
