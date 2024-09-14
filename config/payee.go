package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

// each of pattern, receiver, payment type is a list of either a
// string or a map where the key is the "pattern" and value is a map
// with meta keys

type TransactionMeta struct {
	Location string `yaml:"location"`
	PayeeRaw string `yaml:"payeeRaw"`
	Note     string `yaml:"note"`
}

type PayeePattern struct {
	Value string
	Type  string
	Meta  *TransactionMeta
}

type PayeePatterns []PayeePattern

type Payee struct {
	Name string

	Account string

	// go text/template template string used to generate the payee text
	Template string `yaml:"template"`

	PayeeRaw PayeePatterns `yaml:"payeeRaw"`

	ReceiverAccountNumber PayeePatterns `yaml:"receiverAccountNumber"`

	PaymentType PayeePatterns `yaml:"paymentType"`

	Meta *TransactionMeta `yaml:"meta"`
}

type PayeeConfig struct {
	Payees map[string]*Payee `yaml:"payees"`
}

// Parse the payee object.  As a shortcut, it can have a single string
// value which is interpreted as PayeeRaw pattern
func (p *Payee) UnmarshalYAML(value *yaml.Node) error {
	var pattern string
	if err := value.Decode(&pattern); err == nil {
		p.PayeeRaw = []PayeePattern{{Value: pattern}}
		return nil
	}

	type payee Payee
	if err := value.Decode((*payee)(p)); err != nil {
		return err
	}

	return nil
}

func (pp *PayeePatterns) UnmarshalYAML(value *yaml.Node) error {
	var pattern string
	if err := value.Decode(&pattern); err == nil {
		*pp = []PayeePattern{{Value: pattern}}
		return nil
	}

	type payeePatterns PayeePatterns
	if err := value.Decode((*payeePatterns)(pp)); err != nil {
		return err
	}

	return nil
}

func (pp *PayeePattern) UnmarshalYAML(value *yaml.Node) error {
	var pattern string
	if err := value.Decode(&pattern); err == nil {
		pp.Value = pattern
		return nil
	}

	rawPattern := make(map[string]*TransactionMeta)
	if err := value.Decode(&rawPattern); err == nil {
		for key, value := range rawPattern {
			pp.Value = key
			pp.Meta = value
			return nil
		}
	}

	return fmt.Errorf("Error parsing PayeePattern")
}

func GetUnknownPayee(payeeRaw string) *Payee {
	return &Payee{
		Name:    "Unknown payee ;" + payeeRaw,
		Account: "Unknown:Account",
	}
}
