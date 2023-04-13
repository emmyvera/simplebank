package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// generate a random number between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// generate a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()

}

// generate a random name for owner
func RandomOwner() string {
	return RandomString(6)
}

// generate a random amount of money
func RandomMoney() int64 {
	return RandomInt(500, 5000)
}

// generate a random currency
func RandomCurrency() string {
	currencies := []string{NGN, USD, EUR, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
