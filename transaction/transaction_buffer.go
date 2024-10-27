package transaction

import (
	cfg "bank-to-ledger/config"
)

type TransactionBuffer struct {
	Transactions []Transaction

	// merge or sum, see Config TwinTransactions
	Twin *cfg.TwinTransaction
}

func (tb *TransactionBuffer) Append(trans Transaction) {
	tb.Transactions = append(tb.Transactions, trans)
}

func (tb TransactionBuffer) Match(trans Transaction) bool {
	return tb.Twin != nil &&
		(tb.Length() < tb.Twin.Limit || tb.Twin.Limit == 0) &&
		trans.Match(tb.Twin.Matchers)
}

func (tb TransactionBuffer) Format() string {
	n := len(tb.Transactions)
	if n == 0 {
		return ""
	}

	var mainTransaction Transaction
	if tb.Twin.UseAnchor {
		mainTransaction = tb.Transactions[0]
	} else {
		mainTransaction = tb.Transactions[n-1]
	}

	var bufferTransactions []Transaction

	if n > 1 {
		if tb.Twin.UseAnchor {
			bufferTransactions = tb.Transactions[1:]
		} else {
			bufferTransactions = tb.Transactions[:n-1]
		}

		return mainTransaction.FormatTrans(TransactionBuffer{Transactions: bufferTransactions, Twin: tb.Twin})
	}

	return mainTransaction.FormatTrans(TransactionBuffer{})
}

func (tb TransactionBuffer) IsEmpty() bool {
	return len(tb.Transactions) == 0
}

func (tb TransactionBuffer) Length() int {
	return len(tb.Transactions)
}

func (buffer TransactionBuffer) getAmountSum() float64 {
	var amountBuffer = 0.0
	for _, tr := range buffer.Transactions {
		amountBuffer += tr.AmountReal
	}

	return amountBuffer
}
