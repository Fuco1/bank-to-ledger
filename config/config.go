package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type SymbolMap struct {
	To      string `yaml:"to"`
	InFront bool   `yaml:"inFront"`
}

type TransactionMeta struct {
	Location string `yaml:"location"`
	PayeeRaw string `yaml:"payeeRaw"`
	Note     string `yaml:"note"`
}

type MultiAccountTo struct {
	Matcher   map[string]string `yaml:"matcher"`
	AccountTo string            `yaml:"accountTo"`
}

type Config struct {
	Accounts Account `yaml:"accounts"`

	Payees map[string]*Payee `yaml:"payees"`

	ToPayeeRaw struct {
		Pattern map[string]string `yaml:"pattern"`
	} `yaml:"toPayeeRaw"`

	ToPayee struct {
		PayeeRaw              map[string]string `yaml:"payeeRaw"`
		PaymentType           map[string]string `yaml:"paymentType"`
		ReceiverAccountNumber map[string]string `yaml:"receiverAccountNumber"`
	} `yaml:"toPayee"`

	ToAccountTo struct {
		Payee map[string]string `yaml:"payee"`
		Multi []MultiAccountTo  `yaml:"multi"`
	} `yaml:"toAccountTo"`

	ToMeta struct {
		Payee    map[string]TransactionMeta `yaml:"payee"`
		PayeeRaw map[string]TransactionMeta `yaml:"payeeRaw"`
	} `yaml:"toMeta"`

	PayeeIsTravel []string `yaml:"payeeIsTravel"`

	Currencies struct {
		SymbolMap map[string]SymbolMap `yaml:"symbolMap"`
	} `yaml:"currencies"`

	Banks map[string]Bank `yaml:"banks"`
}

func LoadConfig(fileName string) Config {
	var cfg Config

	yamlFile, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		panic(err)
	}

	for k, v := range cfg.ToPayeeRaw.Pattern {
		_, exists := cfg.ToPayee.PayeeRaw[v]

		if !exists {
			// check if it maps to an account directly
			_, exists = cfg.ToAccountTo.Payee[v]
		}

		if !exists {
			panic(fmt.Sprintf("Regexp-mapped raw-payee %s (from %s) does not exist in payee-raw to payee map ", v, k))
		}
	}

	for name, payee := range cfg.Payees {
		payee.Name = name
	}

	MapPayees(cfg.Accounts, "", cfg.Payees)

	return cfg
}

func (c Config) ValidateConfig() bool {
	return true
}
