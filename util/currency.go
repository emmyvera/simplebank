package util

// Currency that is supported by the bank
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	NGN = "NGN"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, NGN:
		return true
	}

	return false
}
