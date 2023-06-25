package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcedfghijklmnopqrstuvwxz"

// init initializes the random number generator with the current time as the seed.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between the specified minimum and maximum values.
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of the specified length using characters from the alphabet.
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]

		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generates a random owner name.
func RandomOwner() string {
	return RandomString(8)
}

// RandomMoney generates a random amount of money.
func RandomMoney() int64 {
	return RandomInt(0, 10000)
}

// RandomCurrency generates a random currency code (EUR, USD, ILS).
func RandomCurrency() string {
	cur := []string{"EUR", "USD", "ILS"}
	return cur[rand.Intn(len(cur))]
}

// RandomEmail generates a random email using a random string followed by the '@gmail.com' domain.
func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(7))
}
