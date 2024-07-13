package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// each of pattern, receiver, payment type is a list of either a
// string or a map where the key is the "pattern" and value is a map
// with meta keys

type PayeePattern struct {
	Value string
	Meta  map[string]string
}

type Payee struct {
	Name string

	Account string

	PayeeRaw []PayeePattern `yaml:"payeeRaw"`

	ReceiverAccountNumber []PayeePattern `yaml:"receiverAccountNumber"`

	PaymentType []PayeePattern `yaml:"paymentType"`

	Meta map[string]string `yaml:"meta"`
}

type PayeeConfig struct {
	Payees map[string]*Payee `yaml:"payees"`
}

func (pp *PayeePattern) UnmarshalYAML(value *yaml.Node) error {
	var pattern string
	if err := value.Decode(&pattern); err == nil {
		pp.Value = pattern
		return nil
	}

	rawPattern := make(map[string]map[string]string)
	if err := value.Decode(&rawPattern); err == nil {
		for key, value := range rawPattern {
			pp.Value = key
			pp.Meta = value
			return nil
		}
	}

	return fmt.Errorf("Error parsing PayeePattern %s", value)
}

func LoadPayees(fileName string) []Payee {
	var cfg PayeeConfig

	yamlFile, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		panic(err)
	}

	var payees []Payee
	for name, payee := range cfg.Payees {
		payee.Name = name
		payees = append(payees, *payee)
	}

	return payees
}
