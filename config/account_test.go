package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestMapPayee_with_singleElementInArray(t *testing.T) {
	yamlData := `
payees:
  Airbnb: Airbnb

accounts:
  Expenses:
    Hotel:
      Airbnb: [Airbnb]
`

	var config Config
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Error unmarshalling YAML: %v", err)
	}

	MapPayees(config.Accounts, "", config.Payees)

	assert.Equal(t, "Expenses:Hotel:Airbnb", config.Payees["Airbnb"].Account)
}

func TestMapPayee_with_singleStringElement(t *testing.T) {
	yamlData := `
payees:
  Airbnb: Airbnb

accounts:
  Expenses:
    Hotel:
      Airbnb: Airbnb
`

	var config Config
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Error unmarshalling YAML: %v", err)
	}

	MapPayees(config.Accounts, "", config.Payees)

	assert.Equal(t, "Expenses:Hotel:Airbnb", config.Payees["Airbnb"].Account)
}

func TestMapPayee_with_emptyObject(t *testing.T) {
	yamlData := `
payees:
  Airbnb: Airbnb

accounts:
  Expenses:
    Hotel:
      Airbnb:
`

	var config Config
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Error unmarshalling YAML: %v", err)
	}

	MapPayees(config.Accounts, "", config.Payees)

	assert.Equal(t, "Expenses:Hotel:Airbnb", config.Payees["Airbnb"].Account)
}

func TestMapPayee_with_selfSingleAndSubAccounts(t *testing.T) {
	yamlData := `
payees:
  Pharmacy: Pharmacy
  Dentist: Dentist

accounts:
  Expenses:
    Healthcare:
      self: Pharmacy
      Dentist: Dentist
`

	var config Config
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Error unmarshalling YAML: %v", err)
	}

	MapPayees(config.Accounts, "", config.Payees)

	assert.Equal(t, "Expenses:Healthcare", config.Payees["Pharmacy"].Account)
	assert.Equal(t, "Expenses:Healthcare:Dentist", config.Payees["Dentist"].Account)
}

func TestMapPayee_with_selfArrayAndSubAccounts(t *testing.T) {
	yamlData := `
payees:
  Pharmacy: Pharmacy
  Dentist: Dentist

accounts:
  Expenses:
    Healthcare:
      self:
        - Pharmacy
      Dentist: Dentist
`

	var config Config
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Error unmarshalling YAML: %v", err)
	}

	MapPayees(config.Accounts, "", config.Payees)

	assert.Equal(t, "Expenses:Healthcare", config.Payees["Pharmacy"].Account)
	assert.Equal(t, "Expenses:Healthcare:Dentist", config.Payees["Dentist"].Account)
}

func TestMapPayee_implicitPayeeCreated(t *testing.T) {
	yamlData := `
payees:
  Doctor: Doctor

accounts:
  Expenses:
    Restaurant:
      - Old Mill
      - Qerko
`

	var config Config
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Error unmarshalling YAML: %v", err)
	}

	MapPayees(config.Accounts, "", config.Payees)

	assert.Equal(t, "Expenses:Restaurant", config.Payees["Old Mill"].Account)
	assert.Equal(t, "Expenses:Restaurant", config.Payees["Qerko"].Account)
}
