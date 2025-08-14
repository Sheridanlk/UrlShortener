package random

import (
	"math/rand"
	"time"
)

func NewRandomAlias(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	symbols := []rune("QWERTYUIOPASDFGHJKLZXCVBNM" + "qwertyuiopasdfghjklzxcvbnm" + "1234567890")

	output := make([]rune, size)

	for i := range output {
		output[i] = symbols[rnd.Intn(len(symbols))]
	}

	return string(output)
}
