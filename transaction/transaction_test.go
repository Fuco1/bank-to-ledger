package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatAmount_positive_currencyInFront(t *testing.T) {
	result := formatAmount(12.34, CurrencyInfo{
		Sign:      "$",
		IsInFront: true,
	})

	assert.Equal(t, "$12.34", result)
}

func TestFormatAmount_positive_currencyBehind(t *testing.T) {
	result := formatAmount(12.34, CurrencyInfo{
		Sign:      "Kc",
		IsInFront: false,
	})

	assert.Equal(t, "12.34 Kc", result)
}

func TestFormatAmount_negative_currencyInFront(t *testing.T) {
	result := formatAmount(-12.34, CurrencyInfo{
		Sign:      "$",
		IsInFront: true,
	})

	assert.Equal(t, "-$12.34", result)
}

func TestFormatAmount_negative_currencyBehind(t *testing.T) {
	result := formatAmount(-12.34, CurrencyInfo{
		Sign:      "Kc",
		IsInFront: false,
	})

	assert.Equal(t, "-12.34 Kc", result)
}
