package identity

import (
	"math/rand"
	"time"
)

var runeLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(size int) string {
	str := make([]rune, size)
	for i := range str {
		str[i] = runeLetters[rand.Intn(len(runeLetters))]
	}
	return string(str)
}
