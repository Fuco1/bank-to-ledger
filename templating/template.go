package templating

import (
	"bytes"
	"text/template"
)

type TextTemplateBank struct {
	Name           string
	DisplayName    string
	PayeeName      string
	AccountName    string
	FeeAccountName string
	Templates      map[string]string
}

type TextTemplateTransaction struct {
	DateRaw         string
	PayeeRaw        string
	CurrencyRaw     string
	CurrencyAccount string
	PaymentType     string

	Commodity         string
	CommodityPrice    float64
	CommodityQuantity float64

	AmountReal    float64
	AmountAccount float64
	Fee           float64

	ReceiverAccountNumber string

	NoteForMe       string
	NoteForReceiver string
}

type TextTemplatePayee struct {
	Name    string
	Account string
}

type TextTemplateParams struct {
	Bank        TextTemplateBank
	Transaction TextTemplateTransaction
	Payee       TextTemplatePayee
}

func FormatTextTemplate(tmp string, context TextTemplateParams) string {
	tmpl, err := template.New("templating").Parse(tmp)
	var out bytes.Buffer

	err = tmpl.Execute(&out, context)

	if err != nil {
		panic(err)
	}

	return out.String()
}
