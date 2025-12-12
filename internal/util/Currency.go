package util

const (
	USD = "USD"
	EUR = "EUR"
	INR = "INR"
	YEN = "YEN"
	CAD = "CAD"
	BDT = "BDT"
	BRL = "BRL"
	FJD = "FJD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, INR, YEN, CAD, BDT, BRL, FJD:
		return true
	}
	return false
}
