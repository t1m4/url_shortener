package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"url_shortener/configs"
)

func GetStringHash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}
func GetNewLink(length int) string {
	if length < 0 {
		return ""
	}
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		charIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(configs.UrlCharSet))))
		result[i] = configs.UrlCharSet[charIndex.Int64()]
	}
	return string(result)
}
