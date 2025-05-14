package filenamegenerator

import (
	"math/rand"
	"strings"
)

func GenerateFilename() string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	lenght := 6
	var b strings.Builder
	for i := 0; i < lenght; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()
	return str
}
