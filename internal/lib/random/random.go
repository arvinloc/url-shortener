package random

import "math/rand"

var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewRandomString(amount int) string {
	var result string

	for i := 0; i < amount; i++ {
		result += string(letters[rand.Intn(len(letters))])
	}

	return result
}
