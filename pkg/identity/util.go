package identity

import (
	"math/rand"
	"time"
)

// CR(explodes): could be const string
var runeLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

// CR(explodes): random number generators suck for testing.
// i often create a 
// type Randomer interface {
//	Roll() float64
//	Intn(maxExclusive int) int
// }
// that assists with testing.
// its a pain to pass it around but it makes testing ezpz

// CR(explodes): does this need to be exported?
func RandomString(size int) string {
	str := make([]rune, size)
	for i := range str {
		str[i] = runeLetters[rand.Intn(len(runeLetters))]
	}
	return string(str)
}
