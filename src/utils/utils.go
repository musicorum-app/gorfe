package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"unicode/utf8"
)

func FailOnError(err error) {
	if err != nil {
		fmt.Println("An error ocorrured!")
		fmt.Println(err)
	}
}

func Hash(key string) string {
	hash := sha1.New()
	hash.Write([]byte(key))
	bs := hash.Sum(nil)
	return hex.EncodeToString(bs)
}

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}
