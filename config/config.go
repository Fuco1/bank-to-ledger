package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type SymbolMap struct {
	To      string `yaml:"to"`
	InFront bool   `yaml:"inFront"`
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
	} `yaml:"toAccountTo"`

	ToMeta struct {
		Payee    map[string]TransactionMeta `yaml:"payee"`
		PayeeRaw map[string]TransactionMeta `yaml:"payeeRaw"`
	} `yaml:"toMeta"`

	PayeeIsTravel []string `yaml:"payeeIsTravel"`

	Currencies struct {
		SymbolMap map[string]SymbolMap `yaml:"symbolMap"`
	} `yaml:"currencies"`

	Banks map[string]*Bank `yaml:"banks"`
}

func getBankDisplayName(bank Bank) string {
	bankDisplayName := bank.DisplayName
	if bankDisplayName == "" {
		bankDisplayName = bank.PayeeName
	}
	if bankDisplayName == "" {
		bankDisplayName = bank.Name
	}

	if bankDisplayName == "" {
		panic(fmt.Sprintf("Bank %s has no display name", bank.Name))
	}

	return bankDisplayName
}

// Provide default values for column indices
func (ci *ColumnIndices) UnmarshalYAML(value *yaml.Node) error {
	type columnIndices ColumnIndices

	ind := columnIndices{
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

	if err := value.Decode(&ind); err != nil {
		return err
	}

	*ci = ColumnIndices(ind)

	return nil
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

	if len(cfg.Payees) == 0 {
		cfg.Payees = make(map[string]*Payee)
	}

	for name, bank := range cfg.Banks {
		bank.Name = name
		bank.DisplayName = getBankDisplayName(*bank)

		if bank.PayeeName != "" {
			p, exists := cfg.Payees[bank.PayeeName]
			if exists {
				bank.Payee = p
			}
			if p.Account == "" {
				p.Account = bank.AccountName
			}
		}
	}

	MapPayees(cfg.Accounts, "", cfg.Payees)

	return cfg
}

func (c Config) ValidateConfig() bool {
	for _, payee := range c.Payees {
		if payee.Account == "" && payee.AccountTemplate == "" {
			fmt.Fprintf(os.Stderr, "Payee `%s' has no assigned account\n", payee.Name)
		}
	}

	return true
}
