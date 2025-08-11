package transaction

import (
	cfg "bank-to-ledger/config"

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

func TestTransactionMeta_with_pattern_meta_only(t *testing.T) {
	payee := cfg.Payee{
		Name:    "Tiger",
		Account: "",
		PayeeRaw: []cfg.PayeePattern{
			{
				Value: "^tiger.*?",
				Type:  "",
				Meta: &map[string]string{
					"location": "Prague",
				},
			},
		},
		ReceiverAccountNumber: nil,
		PaymentType:           nil,
		Meta:                  nil,
		// Meta: &map[string]string{
		// 	"payeeRaw": "{{ .Transaction.PayeeRaw }}",
		// },
	}

	transaction := Transaction{
		DateRaw:               "",
		PayeeRaw:              "",
		CurrencyRaw:           "",
		CurrencyAccount:       "",
		PaymentType:           "",
		Commodity:             "",
		CommodityPrice:        0,
		CommodityQuantity:     0,
		AmountReal:            0,
		AmountAccount:         0,
		Fee:                   0,
		ReceiverAccountNumber: "",
		NoteForMe:             "",
		NoteForReceiver:       "",
		config: cfg.Config{
			ToMeta: cfg.ToMetaConfig{
				Payee:    map[string]cfg.TransactionMeta{},
				PayeeRaw: map[string]cfg.TransactionMeta{},
			},
		},
		bank: &cfg.Bank{
			Name: "Foo",
		},
		payee:   &payee,
		pattern: &payee.PayeeRaw[0],
	}

	meta := transaction.GetMeta("Tiger")
	assert.NotNil(t, meta)
	assert.Equal(t, "Prague", meta["location"])
}

func TestTransactionMeta_with_payee_meta_only(t *testing.T) {
	payee := cfg.Payee{
		Name:    "Tiger",
		Account: "",
		PayeeRaw: []cfg.PayeePattern{
			{
				Value: "^tiger.*?",
				Type:  "",
				Meta:  &map[string]string{},
			},
		},
		ReceiverAccountNumber: nil,
		PaymentType:           nil,
		Meta: &map[string]string{
			"payeeRaw": "{{ .Transaction.PayeeRaw }}",
		},
	}

	transaction := Transaction{
		DateRaw:               "",
		PayeeRaw:              "TEST",
		CurrencyRaw:           "",
		CurrencyAccount:       "",
		PaymentType:           "",
		Commodity:             "",
		CommodityPrice:        0,
		CommodityQuantity:     0,
		AmountReal:            0,
		AmountAccount:         0,
		Fee:                   0,
		ReceiverAccountNumber: "",
		NoteForMe:             "",
		NoteForReceiver:       "",
		config: cfg.Config{
			ToMeta: cfg.ToMetaConfig{
				Payee:    map[string]cfg.TransactionMeta{},
				PayeeRaw: map[string]cfg.TransactionMeta{},
			},
		},
		bank: &cfg.Bank{
			Name: "Foo",
		},
		payee:   &payee,
		pattern: &payee.PayeeRaw[0],
	}

	meta := transaction.GetMeta("Tiger")
	assert.NotNil(t, meta)
	assert.Equal(t, "TEST", meta["payeeRaw"])
}

func TestTransactionMeta_with_payee_meta_and_pattern_meta(t *testing.T) {
	payee := cfg.Payee{
		Name:    "Tiger",
		Account: "",
		PayeeRaw: []cfg.PayeePattern{
			{
				Value: "^tiger.*?",
				Type:  "",
				Meta: &map[string]string{
					"location": "Prague",
				},
			},
		},
		ReceiverAccountNumber: nil,
		PaymentType:           nil,
		Meta: &map[string]string{
			"payeeRaw": "{{ .Transaction.PayeeRaw }}",
		},
	}

	transaction := Transaction{
		DateRaw:               "",
		PayeeRaw:              "TEST",
		CurrencyRaw:           "",
		CurrencyAccount:       "",
		PaymentType:           "",
		Commodity:             "",
		CommodityPrice:        0,
		CommodityQuantity:     0,
		AmountReal:            0,
		AmountAccount:         0,
		Fee:                   0,
		ReceiverAccountNumber: "",
		NoteForMe:             "",
		NoteForReceiver:       "",
		config: cfg.Config{
			ToMeta: cfg.ToMetaConfig{
				Payee:    map[string]cfg.TransactionMeta{},
				PayeeRaw: map[string]cfg.TransactionMeta{},
			},
		},
		bank: &cfg.Bank{
			Name: "Foo",
		},
		payee:   &payee,
		pattern: &payee.PayeeRaw[0],
	}

	meta := transaction.GetMeta("Tiger")
	assert.NotNil(t, meta)
	assert.Equal(t, "TEST", meta["payeeRaw"])
	assert.Equal(t, "Prague", meta["location"])
}

func TestTransactionMeta_with_payee_meta_and_pattern_overwrite_meta(t *testing.T) {
	payee := cfg.Payee{
		Name:    "Tiger",
		Account: "",
		PayeeRaw: []cfg.PayeePattern{
			{
				Value: "^tiger.*?",
				Type:  "",
				Meta: &map[string]string{
					"location": "Brno",
				},
			},
		},
		ReceiverAccountNumber: nil,
		PaymentType:           nil,
		Meta: &map[string]string{
			"location": "Prague",
		},
	}

	transaction := Transaction{
		DateRaw:               "",
		PayeeRaw:              "TEST",
		CurrencyRaw:           "",
		CurrencyAccount:       "",
		PaymentType:           "",
		Commodity:             "",
		CommodityPrice:        0,
		CommodityQuantity:     0,
		AmountReal:            0,
		AmountAccount:         0,
		Fee:                   0,
		ReceiverAccountNumber: "",
		NoteForMe:             "",
		NoteForReceiver:       "",
		config: cfg.Config{
			ToMeta: cfg.ToMetaConfig{
				Payee:    map[string]cfg.TransactionMeta{},
				PayeeRaw: map[string]cfg.TransactionMeta{},
			},
		},
		bank: &cfg.Bank{
			Name: "Foo",
		},
		payee:   &payee,
		pattern: &payee.PayeeRaw[0],
	}

	meta := transaction.GetMeta("Tiger")
	assert.NotNil(t, meta)
	assert.Equal(t, "Brno", meta["location"])
}

func TestTransactionAccountToAmountFormatting_with_conversion(t *testing.T) {
	payee := cfg.Payee{
		Name:    "Tiger",
		Account: "",
		PayeeRaw: []cfg.PayeePattern{
			{
				Value: "^tiger.*?",
				Type:  "",
				Meta: &map[string]string{
					"location": "Prague",
				},
			},
		},
		ReceiverAccountNumber: nil,
		PaymentType:           nil,
		Meta:                  nil,
	}

	transaction := Transaction{
		DateRaw:               "",
		PayeeRaw:              "",
		CurrencyRaw:           "PLN",
		CurrencyAccount:       "CZK",
		PaymentType:           "",
		Commodity:             "",
		CommodityPrice:        0,
		CommodityQuantity:     0,
		AmountReal:            -10,
		AmountAccount:         -250,
		Fee:                   0,
		ReceiverAccountNumber: "",
		NoteForMe:             "",
		NoteForReceiver:       "",
		config: cfg.Config{
			ToMeta: cfg.ToMetaConfig{
				Payee:    map[string]cfg.TransactionMeta{},
				PayeeRaw: map[string]cfg.TransactionMeta{},
			},
		},
		bank: &cfg.Bank{
			Name: "Foo",
		},
		payee:   &payee,
		pattern: &payee.PayeeRaw[0],
	}

	meta := transaction.GetMeta("Tiger")
	assert.NotNil(t, meta)
	assert.Equal(t, "10.00 PLN @@ 250.00 Kc", transaction.FormatAmountRealInverted(nil))
}
