package config

import (
	"fmt"
	"os"
)

type Account map[string]interface{}

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

func assignAccount(payeeName string, accountName string, payees map[string]*Payee) {
	val, exists := payees[payeeName]

	if exists {
		if val.Account != "" {
			fmt.Fprintf(os.Stderr, "Payee `%s' already has assigned account `%s', trying to assign `%s'\n", payeeName, val.Account, accountName)
		} else {
			val.Account = accountName
		}
	} else {
		// If we simply assign a "payee" to an account without
		// explicit payee definition, we create an implicit definition
		// here.  This way for simple payees we only need to add one
		// line ("category" or the account) to make it recognized it
		// in the output
		payees[payeeName] = &Payee{
			Name:     payeeName,
			Account:  accountName,
			PayeeRaw: []PayeePattern{{Value: "^" + payeeName + "$"}},
		}
	}
}

// recursively travers accounts hierarchy and map payees in the leaves to the account path
func MapPayees(account Account, path string, payees map[string]*Payee) {
	for key, value := range account {
		var delim string
		if path == "" {
			delim = ""
		} else {
			delim = ":"
		}

		accountName := path + delim + key
		if key == "self" {
			accountName = path
		}

		switch value.(type) {
		case nil:
			assignAccount(key, accountName, payees)
		case Account:
			MapPayees(value.(Account), accountName, payees)
		case string:
			payee := value.(string)
			assignAccount(payee, accountName, payees)
		case []interface{}:
			p, _ := convertToStrings(value.([]interface{}))
			for _, payee := range p {
				assignAccount(payee, accountName, payees)
			}
		}
	}
}
