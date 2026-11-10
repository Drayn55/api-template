package currencyutil

const BaseCurrency = "IDR"

func EffectiveRate(currencyCode string, currencyRate float64) float64 {
	if currencyRate > 0 {
		return currencyRate
	}
	return 1
}

func ToBase(amount float64, currencyCode string, currencyRate float64) float64 {
	return amount * EffectiveRate(currencyCode, currencyRate)
}

func ToBasePtr(amount *float64, currencyCode string, currencyRate float64) *float64 {
	if amount == nil {
		return nil
	}
	v := ToBase(*amount, currencyCode, currencyRate)
	return &v
}
