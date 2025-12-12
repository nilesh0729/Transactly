package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"

 //Random int generator between min and max
func RandomInt(min, max int64) int64 {
	return (min + rand.Int63n(max-min+1))
}

// Random string generator generates a string of "n" length
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i:=0; i<n; i++ {
		c:= alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Generating random Owner
func RandomOwner()string{
	return RandomString(6)
}

// Generate a Random Balance
func RandomBalance()int64{
	return RandomInt(0,1000)
}

//Generate a Random Currency
func RandomCurrency()string{ 
    Currencies := []string{"INR","CAD","EUR","USD","YEN"}
	N:= len(Currencies)

	return Currencies[rand.Intn(N)]
}

func RandomEmail() string{
	return fmt.Sprintf("%s@gmail.com",RandomString(8))
}
