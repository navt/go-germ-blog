package utility

import (
	"crypto/rand"
	"fmt"
	"time"
)

// GenerateId Генерация строки содержащей число
func GenerateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
func IdFromTime() string {
	t := time.Now()
	return t.Format("2006-01-02_15-04-05")
}
