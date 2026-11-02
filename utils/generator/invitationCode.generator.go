package generator

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateInviteCode(slug string) string {

	runes := []rune(slug)

	totalChars := len(runes)
	n := totalChars / 2

	prefixString := string(runes[:n])
	suffixString := string(runes[n:])

	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return prefixString+suffixString+hex.EncodeToString(b)
}