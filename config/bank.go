package config

import (
	"log"
)

type ColumnIndices struct {
	DateRaw               int `yaml:"dateRaw"`
	PayeeRaw              int `yaml:"payeeRaw"`
	CurrencyRaw           int `yaml:"currencyRaw"`
	CurrencyAccount       int `yaml:"currencyAccount"`
	PaymentType           int `yaml:"paymentType"`
	Commodity             int `yaml:"commodity"`
	CommodityPrice        int `yaml:"commodityPrice"`
	CommodityQuantity     int `yaml:"commodityQuantity"`
	AmountReal            int `yaml:"amountReal"`
	AmountAccount         int `yaml:"amountAccount"`
	Fee                   int `yaml:"fee"`
	ReceiverAccountNumber int `yaml:"receiverAccountNumber"`
	NoteForMe             int `yaml:"noteForMe"`
	NoteForReceiver       int `yaml:"noteForReceiver"`
}

type ColumnNames struct {
	DateRaw               string `yaml:"dateRaw"`
	PayeeRaw              string `yaml:"payeeRaw"`
	CurrencyRaw           string `yaml:"currencyRaw"`
	CurrencyAccount       string `yaml:"currencyAccount"`
	PaymentType           string `yaml:"paymentType"`
	Commodity             string `yaml:"commodity"`
	CommodityPrice        string `yaml:"commodityPrice"`
	CommodityQuantity     string `yaml:"commodityQuantity"`
	AmountReal            string `yaml:"amountReal"`
	AmountAccount         string `yaml:"amountAccount"`
	Fee                   string `yaml:"fee"`
	ReceiverAccountNumber string `yaml:"receiverAccountNumber"`
	NoteForMe             string `yaml:"noteForMe"`
	NoteForReceiver       string `yaml:"noteForReceiver"`
}

type Matcher struct {
	DateRaw               string `yaml:"dateRaw"`
	Payee                 string `yaml:"payee"`
	PayeeRaw              string `yaml:"payeeRaw"`
	CurrencyRaw           string `yaml:"currencyRaw"`
	CurrencyAccount       string `yaml:"currencyAccount"`
	PaymentType           string `yaml:"paymentType"`
	AmountReal            string `yaml:"amountReal"`
	AmountAccount         string `yaml:"amountAccount"`
	Fee                   string `yaml:"fee"`
	ReceiverAccountNumber string `yaml:"receiverAccountNumber"`
	NoteForMe             string `yaml:"noteForMe"`
	NoteForReceiver       string `yaml:"noteForReceiver"`
}

type TwinTransaction struct {
	// type can be
	// - `sum` for adding the amount to previous transaction's amount
	//   (produce 2 line transaction only)
	// - `merge` for adding the account and amount to previous
	//   transaction (produces transaction with multiple lines)
	Type string `yaml:"type"`

	Inverted bool `yaml:"inverted"`

	Anchor []Matcher `yaml:"anchor"`

	Matchers []Matcher `yaml:"matchers"`

	UseAnchor bool `yaml:"useAnchor"`

	Limit int `yaml:"limit"`
}

type IgnoredTransactions struct {
	Matchers []Matcher `yaml:"matchers"`
}

type Bank struct {
	// Name of the bank (config key)
	Name string

	// Display name of the bank.  Defaults to PayeeName and then Name.
	DisplayName string `yaml:"displayName"`

	// Payee name of this bank's payee.  If transaction matches this
	// payee, we know the transaction is between user's own accounts
	// and we only record one side (because presumably it would appear
	// in both exports).
	PayeeName string `yaml:"payee"`

	// Resolved full Payee object, based on the PayeeName
	Payee *Payee `yaml:"-"`

	// Name of the checking account representing at this bank
	AccountName string `yaml:"accountName"`

	// Name of the account accruing this bank's fees
	FeeAccountName string `yaml:"feeAccountName"`

	Templates map[string]string `yaml:"templates"`

	// Pattern used to parse date from DateRaw column
	DatePatternFrom string `yaml:"datePatternFrom"`

	// Columns used to auto-identify the bank from a csv file
	IdentifyingColumns []string `yaml:"identifyingColumns"`

	FileNamePattern string `yaml:"fileNamePattern"`

	ColumnNames ColumnNames `yaml:"columnNames"`

	ColumnIndices ColumnIndices `yaml:"columnIndices"`

	TwinTransactions []TwinTransaction `yaml:"twinTransactions"`

	IgnoredTransactions []IgnoredTransactions `yaml:"ignoredTransactions"`
}

func (b Bank) NamesToIndices(header []string) ColumnIndices {
	if (b.ColumnIndices != ColumnIndices{}) {
		return b.ColumnIndices
	}

	indices := ColumnIndices{
		DateRaw:               -1,
		PayeeRaw:              -1,
		CurrencyRaw:           -1,
		CurrencyAccount:       -1,
		PaymentType:           -1,
		Commodity:             -1,
		CommodityPrice:        -1,
		CommodityQuantity:     -1,
		AmountReal:            -1,
		AmountAccount:         -1,
		Fee:                   -1,
		ReceiverAccountNumber: -1,
		NoteForMe:             -1,
		NoteForReceiver:       -1,
	}

	for i, v := range header {
		if v == b.ColumnNames.DateRaw {
			indices.DateRaw = i
		}
		if v == b.ColumnNames.PayeeRaw {
			indices.PayeeRaw = i
		}
		if v == b.ColumnNames.CurrencyRaw {
			indices.CurrencyRaw = i
		}
		if v == b.ColumnNames.CurrencyAccount {
			indices.CurrencyAccount = i
		}
		if v == b.ColumnNames.PaymentType {
			indices.PaymentType = i
		}
		if v == b.ColumnNames.Commodity {
			indices.Commodity = i
		}
		if v == b.ColumnNames.CommodityPrice {
			indices.CommodityPrice = i
		}
		if v == b.ColumnNames.CommodityQuantity {
			indices.CommodityQuantity = i
		}
		if v == b.ColumnNames.AmountReal {
			indices.AmountReal = i
		}
		if v == b.ColumnNames.AmountAccount {
			indices.AmountAccount = i
		}
		if v == b.ColumnNames.Fee {
			indices.Fee = i
		}
		if v == b.ColumnNames.ReceiverAccountNumber {
			indices.ReceiverAccountNumber = i
		}
		if v == b.ColumnNames.NoteForMe {
			indices.NoteForMe = i
		}
		if v == b.ColumnNames.NoteForReceiver {
			indices.NoteForReceiver = i
		}
	}

	return indices
}

// check if leading columns in order correspond to identifying columns
func GetBankConfig(header []string, banks map[string]*Bank) (*Bank, bool) {
	for name, bank := range banks {
		if len(bank.IdentifyingColumns) == 0 {
			continue
		}

		if len(bank.IdentifyingColumns) > len(header) {
			continue
		}

		identifyingColumnsMatch := true
		for i, v := range bank.IdentifyingColumns {
			if header[i] != v {
				identifyingColumnsMatch = false
				break
			}
		}

		if identifyingColumnsMatch {
			// TODO: should never be empty
			if bank.Name == "" {
				bank.Name = name
			}
			return bank, true
		}
	}

	return &Bank{}, false
}

func (b Bank) ValidateBankConfig() bool {
	if b.DatePatternFrom == "" {
		log.Fatalf("DatePatternFrom is not set for bank %s", b.Name)
	}

	return true
}
