package fraas_helpers

import (
	"math/rand"
	"time"
)

const dictionary = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func GeneratePassword(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = dictionary[rand.Intn(len(dictionary))]
	}
	return string(b)
}

func GenerateSecret(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = dictionary[rand.Intn(len(dictionary)-10)]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
