package model

import (
	"math/rand"
	"time"
)

const (
	// 62 = 26 + 26 + 10
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 6 bits to represent 62 index
	letterIdxBits = 6
	// All 1-bits, as many as letterBits
	letterIdxMask = 1<<letterIdxBits - 1

	IDSize = 10
)

func generateRandomString(n int) string {
	src := rand.NewSource(time.Now().UnixNano() + int64(rand.Intn(1024)))
	b := make([]byte, n)
	l := len(letters)
	for i := 0; i < n; {
		idx := int(src.Int63() & letterIdxMask)
		if idx < l {
			b[i] = letters[idx]
			i++
		}
	}
	return string(b)
}

func generateID() string {
	return generateRandomString(IDSize)
}
