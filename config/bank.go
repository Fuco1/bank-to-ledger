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
	AmountReal            string `yaml:"amountReal"`
	AmountAccount         string `yaml:"amountAccount"`
	Fee                   string `yaml:"fee"`
	ReceiverAccountNumber string `yaml:"receiverAccountNumber"`
	NoteForMe             string `yaml:"noteForMe"`
	NoteForReceiver       string `yaml:"noteForReceiver"`
}

type Matcher struct {
	Column string `yaml:"column"`
	Value  string `yaml:"value"`
}

type TwinTransactions struct {
	// type can be
	// - `sum` for adding the amount to previous transaction's amount
	//   (produce 2 line transaction only)
	// - `merge` for adding the account and amount to previous
	//   transaction (produces transaction with multiple lines)
	Type string `yaml:"type"`

	Inverted bool `yaml:"inverted"`

	Matchers []Matcher `yaml:"matchers"`
}

type IgnoredTransactions struct {
	Matchers []Matcher `yaml:"matchers"`
}

type Bank struct {
	// Name of the bank
	Name string

	// Name of the checking account representing at this bank
	Account string `yaml:"checkingAccountName"`

	// Name of the account accruing this bank's fees
	FeeAccount string `yaml:"feeAccountName"`

	Templates map[string]string `yaml:"templates"`

	// Pattern used to parse date from DateRaw column
	DatePatternFrom string `yaml:"datePatternFrom"`

	// Columns used to auto-identify the bank from a csv file
	IdentifyingColumns []string `yaml:"identifyingColumns"`

	ColumnNames ColumnNames `yaml:"columnNames"`

	ColumnIndices ColumnIndices `yaml:"columnIndices"`

	TwinTransactions []TwinTransactions `yaml:"twinTransactions"`

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
func GetBankConfig(header []string, banks map[string]Bank) (Bank, bool) {
	for name, bank := range banks {
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
			if bank.Name == "" {
				bank.Name = name
			}
			return bank, true
		}
	}

	return Bank{}, false
}

func (b Bank) ValidateBankConfig() bool {
	if b.DatePatternFrom == "" {
		log.Fatalf("DatePatternFrom is not set for bank %s", b.Name)
	}

	return true
}
