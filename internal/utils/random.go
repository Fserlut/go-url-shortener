package random

import (
	"fmt"
	"math/rand"
	"time"
)

func GetShortURL() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 8+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2:8]
}
