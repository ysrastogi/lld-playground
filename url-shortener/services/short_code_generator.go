package services

import (
	"math/big"
	"crypto/rand"
)

type ShortCodeStrategy interface {
	generate(shortCode string) string
}

type RandomShortCode struct {
	longURL string
	userId string
	length int
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (r *RandomShortCode) generate() string {
	if r.length <= 0 {
		r.length = 6
	}

	result := make([]byte, r.length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < r.length; i++ {
		num, _ := rand.Int(rand.Reader, charsetLen)
		result[i] = charset[num.Int64()]
	}

	return string(result)
}
