package random

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomSecretToken() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:32]
}
