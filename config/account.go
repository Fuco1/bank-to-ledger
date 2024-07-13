package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Account map[string]interface{}

type AccountsConfig struct {
	Accounts Account `yaml:"accounts"`
}

func convertToStrings(input []interface{}) ([]string, error) {
	var result []string
	for _, v := range input {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("non-string element found: %v", v)
		}
		result = append(result, str)
	}
	return result, nil
}

// recursively travers accounts hierarchy and map payees in the leaves to the account path
func MapPayees(account Account, path string, payees map[string]*Payee) {
	for key, value := range account {
		switch value.(type) {
		case Account:
			MapPayees(value.(Account), path+key+":", payees)
		case []interface{}:
			p, _ := convertToStrings(value.([]interface{}))
			for _, payee := range p {
				val, exists := payees[payee]
				if exists {
					if val.Account != "" {
						fmt.Fprintf(os.Stderr, "Payee `%s' already has assigned account `%s', trying to assign `%s'\n", payee, val.Account, path+key)
					} else {
						val.Account = path + key
					}
				} else {
					fmt.Fprintf(os.Stderr, "Payee `%s' not defined\n", payee)
				}
			}
		}
	}
}

func LoadAccounts(fileName string) AccountsConfig {
	var cfg AccountsConfig

	yamlFile, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
