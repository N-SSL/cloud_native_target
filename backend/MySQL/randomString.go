package MySQL

import (
	"math/rand"
	"time"
)

var DefaultLetters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

// RandomString returns a random string with a fixed length
func RandomString(n int) string {
	var letters []rune

	rand.Seed(time.Now().UnixNano())

	letters = DefaultLetters

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
