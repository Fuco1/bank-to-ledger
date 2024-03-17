package transaction

import (
	cfg "bank-to-ledger/config"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
	"text/template"
)

type FormattableTransaction interface {
	FormatTrans(buffer TransactionBuffer) string
	IsTwinTransaction() string
}

type Transaction struct {
	DateRaw         string
	PayeeRaw        string
	CurrencyRaw     string
	CurrencyAccount string
	PaymentType     string

	AmountReal    float64
	AmountAccount float64
	Fee           float64

	NoteForMe       string
	NoteForReceiver string

	ReceiverAccountNumber string

	Meta cfg.TransactionMeta

	config cfg.Config
	bank   cfg.Bank
}

type TransactionBuffer struct {
	Transactions []Transaction

	// merge or sum, see Config TwinTransactions
	Twin cfg.TwinTransactions
}

type CurrencyInfo struct {
	Sign      string
	IsInFront bool
}

func normalizeAmount(amount string) string {
	amountNoComma := strings.ReplaceAll(amount, ",", ".")
	return strings.ReplaceAll(amountNoComma, " ", "")
}

func FromCsvRecord(record []string, config cfg.Config, bank cfg.Bank) Transaction {
	ci := bank.ColumnIndices

	amountAccountNormalized := normalizeAmount(record[ci.AmountAccount])
	amountRealNormalized := normalizeAmount(record[ci.AmountReal])

	if amountRealNormalized == "" {
		amountRealNormalized = amountAccountNormalized
	}

	AmountAccount, _ := strconv.ParseFloat(amountAccountNormalized, 64)
	AmountReal, _ := strconv.ParseFloat(amountRealNormalized, 64)

	var Fee float64
	if ci.Fee == -1 || record[ci.Fee] == "" {
		Fee = 0
	} else {
		Fee, _ = strconv.ParseFloat(record[ci.Fee], 64)
	}

	currencyRaw := record[ci.CurrencyRaw]

	currencyAccount := ""
	if ci.CurrencyAccount != -1 {
		currencyAccount = record[ci.CurrencyAccount]
	}
	if currencyAccount == "" {
		currencyAccount = currencyRaw
	}

	noteForMe := ""
	if ci.NoteForMe != -1 {
		noteForMe = record[ci.NoteForMe]
	}

	noteForReceiver := ""
	if ci.NoteForReceiver != -1 {
		noteForReceiver = record[ci.NoteForReceiver]
	}

	return Transaction{
		DateRaw:         record[ci.DateRaw],
		PaymentType:     record[ci.PaymentType],
		CurrencyRaw:     currencyRaw,
		CurrencyAccount: currencyAccount,
		PayeeRaw:        record[ci.PayeeRaw],

		AmountAccount: AmountAccount,
		AmountReal:    AmountReal,
		Fee:           Fee,

		NoteForMe:       noteForMe,
		NoteForReceiver: noteForReceiver,

		ReceiverAccountNumber: record[ci.ReceiverAccountNumber],

		config: config,
		bank:   bank,
	}
}

func (t Transaction) FormatDate() string {
	tt, _ := time.Parse(t.bank.DatePatternFrom, t.DateRaw)
	return tt.Format("2006/01/02")
}

func (t Transaction) GetCurrency() CurrencyInfo {
	return t.GetCurrencyBySymbol(t.CurrencyRaw)
}

func (t Transaction) GetCurrencyBySymbol(currency string) CurrencyInfo {
	if currency == "" {
		currency = t.CurrencyAccount
	}

	info := CurrencyInfo{currency, false}

	symbolMap, exists := t.config.Currencies.SymbolMap[currency]
	if exists {
		if info.Sign != "" {
			info.Sign = symbolMap.To
		}
		info.IsInFront = symbolMap.InFront
	}

	return info
}

func (t Transaction) GetExchangeRate() float64 {
	return t.AmountAccount / t.AmountReal
}

func formatAmount(amount float64, ci CurrencyInfo) string {
	if ci.IsInFront {
		return fmt.Sprintf("%s%.2f", ci.Sign, amount)
	} else {
		return fmt.Sprintf("%.2f %s", amount, ci.Sign)
	}
}

func formatExchangeRate(exchangeRate float64, ci CurrencyInfo) string {
	if ci.Sign != "Kc" {
		return fmt.Sprintf(" @ %.6f Kc", exchangeRate)
	} else {
		return ""
	}
}

func (t Transaction) formatAmountReal(amount float64) string {
	ci := t.GetCurrency()
	return t.formatAmountRealWithCurrency(amount, ci)
}

func (t Transaction) formatAmountRealWithCurrency(amount float64, ci CurrencyInfo) string {
	return formatAmount(amount, ci) + formatExchangeRate(t.GetExchangeRate(), ci)
}

func (t Transaction) FormatAmountReal() string {
	return t.formatAmountReal(t.AmountReal)
}

func (t Transaction) FormatFee() string {
	if t.Fee != 0 {
		return t.formatAmountRealWithCurrency(
			-t.Fee,
			t.GetCurrencyBySymbol(t.CurrencyAccount),
		)
	} else {
		return ""
	}
}

func (t Transaction) FormatAmountRealInverted(buffer *TransactionBuffer) string {
	amount := -t.AmountReal

	if buffer != nil && buffer.Twin.Type == "sum" {
		var amountBuffer = 0.0
		for _, tr := range buffer.Transactions {
			amountBuffer -= tr.AmountReal
		}
		amount = amount + amountBuffer
	}

	return t.formatAmountReal(amount)
}

func (t Transaction) GetPayee() (string, bool) {
	payee, exists := t.config.ToPayee.PaymentType[t.PaymentType]

	if !exists {
		payee, exists = t.config.ToPayee.ReceiverAccountNumber[t.ReceiverAccountNumber]
	}

	if !exists {
		payee, exists = t.config.ToPayee.PayeeRaw[t.PayeeRaw]
	}

	if !exists {
		for pattern, payeeRaw := range t.config.ToPayeeRaw.Pattern {
			match, _ := regexp.MatchString("(?i)"+pattern, t.PayeeRaw)
			if match {
				payee, exists = t.config.ToPayee.PayeeRaw[payeeRaw]

				if !exists {
					_, exists = t.config.ToAccountTo.Payee[payeeRaw]
					if exists {
						payee = payeeRaw
					}
				}
				break
			}
		}
	}

	// check if this raw payee is mapped to an account.  If it is, we
	// assume that the raw payee and final payee are the same and no
	// processing is necessary.
	if !exists {
		_, exists = t.config.ToAccountTo.Payee[t.PayeeRaw]
		if exists {
			payee = t.PayeeRaw
		}
	}

	if !exists {
		payee = "Unknown payee ;" + t.PayeeRaw
	}

	return payee, exists
}

func (t Transaction) GetNote() string {
	payee, _ := t.GetPayee()
	note := []string{"(^.^)"}

	if payee == "RegioJet" {
		note = append(note, "check if credit")
	}

	if slices.Contains(t.config.PayeeIsTravel, payee) {
		note = append(note, "add to/from/location")
	}

	if slices.Contains([]string{"Alza", "Tesco", "Tiger"}, payee) {
		note = append(note, "what did you get from there?")
	}

	if payee == "Unknown hotel" {
		note = append(note, "add which hotel it is")
	}

	if payee == "Small chinese shop" {
		note = append(
			note,
			"this is prob the small shop next to your place but check",
		)
	}

	if len(note) > 1 {
		return " " + strings.Join(note, ", ")
	} else {
		return ""
	}
}

func (t Transaction) resolveTemplate(template string) string {
	val, exists := t.bank.Templates[template[1:len(template)-1]]

	if exists {
		return val
	}

	return template
}

func (t Transaction) GetAccountTo() string {
	payee, _ := t.GetPayee()
	acc, exists := t.config.ToAccountTo.Payee[payee]

	// Try the multi-matchers here
	if !exists {
		for _, mmatcher := range t.config.ToAccountTo.Multi {
			matchers := make([]cfg.Matcher, 1)
			for k, v := range mmatcher.Matcher {
				matchers = append(matchers, cfg.Matcher{
					Column: k,
					Value:  v,
				})
			}
			if t.Match(matchers) {
				acc = mmatcher.AccountTo
				exists = true
				break
			}
		}
	}

	if !exists {
		return "Unknown:Account"
	}

	return t.resolveTemplate(acc)
}

func (t Transaction) GetAccountFrom() string {
	account := t.bank.Account
	if account == "" {
		account = "Unknown:AccountFrom"
	}

	return account
}

func (t Transaction) Match(matchers []cfg.Matcher) bool {
	isMatch := true

	for _, matcher := range matchers {
		switch matcher.Column {
		case "PaymentType":
			isMatch = isMatch && t.PaymentType == matcher.Value
		case "ReceiverAccountNumber":
			isMatch = isMatch && t.ReceiverAccountNumber == matcher.Value
		case "PayeeRaw":
			isMatch = isMatch && t.PayeeRaw == matcher.Value
		case "NoteForMe":
			isMatch = isMatch && t.NoteForMe == matcher.Value
		case "NoteForReceiver":
			isMatch = isMatch && t.NoteForReceiver == matcher.Value
		default:
			return false
		}
	}

	return isMatch
}

func (t Transaction) IsTwinTransaction() *cfg.TwinTransactions {
	ttConfig := t.bank.TwinTransactions
	for _, tt := range ttConfig {
		if len(tt.Matchers) == 0 {
			continue
		}

		if t.Match(tt.Matchers) {
			return &tt
		}
	}

	return nil
}

func (t Transaction) IsIgnored() bool {
	ignored := t.bank.IgnoredTransactions
	for _, ignoreDef := range ignored {
		if len(ignoreDef.Matchers) == 0 {
			continue
		}

		if t.Match(ignoreDef.Matchers) {
			return true
		}
	}

	return false
}

func (tb *TransactionBuffer) Append(trans Transaction) {
	tb.Transactions = append(tb.Transactions, trans)
}

func (t Transaction) FormatTwinTransaction(buffer TransactionBuffer) string {
	if buffer.Twin.Type == "merge" {
		lines := make([]string, 1)
		for _, tt := range buffer.Transactions {
			var amount string
			if buffer.Twin.Inverted {
				amount = tt.FormatAmountRealInverted(nil)
			} else {
				amount = tt.FormatAmountReal()
			}

			lines = append(lines, fmt.Sprintf(
				"    %s  %s",
				tt.GetAccountTo(),
				amount,
			))
		}
		return strings.Join(lines, "\n")
	}

	return ""
}

func (t Transaction) getMetaFromStruct(
	meta cfg.TransactionMeta,
	metaOut map[string]string,
) {
	if meta.Location != "" {
		metaOut["Location"] = meta.Location
	}

	if meta.PayeeRaw != "" {
		if meta.PayeeRaw == "%payeeRaw%" {
			metaOut["PayeeRaw"] = t.PayeeRaw
		} else {
			metaOut["PayeeRaw"] = meta.PayeeRaw
		}
	}

	if meta.Note != "" {
		metaOut["Note"] = meta.Note
	}
}

func (t Transaction) GetMeta(payee string) map[string]string {
	metaOut := make(map[string]string)

	meta, exists := t.config.ToMeta.Payee[payee]
	if exists {
		t.getMetaFromStruct(meta, metaOut)
	}

	meta, exists = t.config.ToMeta.PayeeRaw[t.PayeeRaw]
	if exists {
		t.getMetaFromStruct(meta, metaOut)
	}

	return metaOut
}

type TemplateContext struct {
	Transaction     Transaction
	Date            string
	Payee           string
	Note            string
	Meta            string
	AccountTo       string
	AccountToAmount string
	FeeAmount       string
	AccountFee      string
	AmountTotal     string
	TwinTransaction string
	AccountFrom     string
}

func (t Transaction) FormatTrans(buffer TransactionBuffer) string {
	payee, _ := t.GetPayee()
	meta := t.GetMeta(payee)

	metaLines := make([]string, 1)
	for k, v := range meta {
		metaLines = append(metaLines, fmt.Sprintf("    ; %s: %s\n", k, v))
	}

	tmpl, err := template.New("transaction").Parse(`{{ .Date }} * {{ .Payee }}{{ .Note }}
{{ .Meta }}    {{ .AccountTo }}  {{ .AccountToAmount }}
{{- if .FeeAmount }}
    {{ .AccountFee }}  {{ .FeeAmount }}
{{- end }}
{{- if .TwinTransaction -}}
{{ .TwinTransaction }}
{{- end }}
    {{ .AccountFrom }}{{ if or .FeeAmount (ne .Transaction.CurrencyRaw .Transaction.CurrencyAccount) }}  {{ .AmountTotal }}{{ end }}
`)

	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	context := TemplateContext{
		t,
		t.FormatDate(),
		payee,
		t.GetNote(),
		strings.Join(metaLines, ""),
		t.GetAccountTo(),
		t.FormatAmountRealInverted(&buffer),
		t.FormatFee(),
		t.bank.FeeAccount,
		t.formatAmountRealWithCurrency(t.AmountAccount+t.Fee, t.GetCurrencyBySymbol(t.CurrencyAccount)),
		t.FormatTwinTransaction(buffer),
		t.GetAccountFrom(),
	}

	err = tmpl.Execute(&out, context)

	if err != nil {
		panic(err)
	}

	return out.String()

	// return fmt.Sprintf(
	// 	"%s * %s%s\n%s    %s  %s\n%s    %s\n",
	// 	t.FormatDate(),
	// 	payee,
	// 	t.GetNote(),
	// 	strings.Join(metaLines, ""),
	// 	t.GetAccountTo(),
	// 	t.FormatAmountRealInverted(&buffer),
	// 	t.FormatTwinTransaction(buffer),
	// 	t.GetAccountFrom(),
	// )
}
