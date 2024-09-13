package config

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func getExpected() Config {
	return Config{
		Payees: map[string]*Payee{
			"TIGER": {
				Name:    "",
				Account: "",
				PayeeRaw: []PayeePattern{
					{
						Value: "^tiger.*?",
						Type:  "",
						Meta:  nil,
					},
				},
				ReceiverAccountNumber: nil,
				PaymentType:           nil,
				Meta:                  nil,
			},
		},
	}
}

func TestUnmarshalPayee_implicit_payeeRaw(t *testing.T) {
	yamlData := `
payees:
  TIGER: '^tiger.*?'
`

	var config Config
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Error unmarshalling YAML: %v", err)
	}

	expected := getExpected()
	if !reflect.DeepEqual(config.Payees, expected.Payees) {
		t.Errorf("Unmarshalled config does not match expected config. Got %+v, expected %+v", config.Payees, expected.Payees)
	}
}

func TestUnmarshalPayee_payeeRaw_string(t *testing.T) {
	yamlData := `
payees:
  TIGER:
    payeeRaw: '^tiger.*?'
`

	var config Config
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Error unmarshalling YAML: %v", err)
	}

	expected := getExpected()
	if !reflect.DeepEqual(config.Payees, expected.Payees) {
		t.Errorf("Unmarshalled config does not match expected config. Got %+v, expected %+v", config.Payees, expected.Payees)
	}
}

func TestUnmarshalPayee_receiverAccountNumber_string(t *testing.T) {
	yamlData := `
payees:
  TIGER:
    receiverAccountNumber: '1234/567'
`

	var config Config
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Error unmarshalling YAML: %v", err)
	}

	expected := getExpected()
	expected.Payees["TIGER"].ReceiverAccountNumber = PayeePatterns{{Value: "1234/567"}}
	expected.Payees["TIGER"].PayeeRaw = nil

	if !reflect.DeepEqual(config.Payees, expected.Payees) {
		t.Errorf("Unmarshalled config does not match expected config. Got %+v, expected %+v", config.Payees, expected.Payees)
	}
}

func TestUnmarshalPayee_payeeRaw_array(t *testing.T) {
	yamlData := `
payees:
  TIGER:
    payeeRaw: ['^tiger.*?']
`

	var config Config
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Error unmarshalling YAML: %v", err)
	}

	expected := getExpected()
	if !reflect.DeepEqual(config.Payees, expected.Payees) {
		t.Errorf("Unmarshalled config does not match expected config. Got %+v, expected %+v", config.Payees, expected.Payees)
	}
}

func TestUnmarshalPayee_payeeRaw_with_meta(t *testing.T) {
	yamlData := `
payees:
  TIGER:
    payeeRaw:
      - '^tiger.*?':
          location: Prague
`

	var config Config
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Error unmarshalling YAML: %v", err)
	}

	var expected Config = Config{
		Payees: map[string]*Payee{
			"TIGER": {
				Name:    "",
				Account: "",
				PayeeRaw: []PayeePattern{
					{
						Value: "^tiger.*?",
						Type:  "",
						Meta: &TransactionMeta{
							Location: "Prague",
						},
					},
				},
				ReceiverAccountNumber: nil,
				PaymentType:           nil,
				Meta:                  nil,
			},
		},
	}

	if !reflect.DeepEqual(config.Payees, expected.Payees) {
		t.Errorf("Unmarshalled config does not match expected config. Got %+v, expected %+v", config.Payees, expected.Payees)
	}
}
