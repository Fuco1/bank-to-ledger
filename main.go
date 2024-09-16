package main

import (
	cfg "bank-to-ledger/config"
	t "bank-to-ledger/transaction"
	"encoding/csv"
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"path/filepath"
	"regexp"
	s "strings"
	//	"github.com/sanity-io/litter"
)

type Options struct {
	Verbose []bool `long:"verbose" description:"Verbose output"`

	Config string `long:"config" description:"Transactions config" default:"config.yaml"`

	HasHeader bool `long:"has-header" description:"Whether first line of csv is header"`

	HasNoHeader bool `long:"has-no-header" description:"Whether first line of csv is header"`

	BankName string `long:"bank-name" description:"Bank name used to determine csv format."`
}

func readCsv(fileName string, options Options, config cfg.Config) ([]t.Transaction, *cfg.Bank) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Parsing error `%s`, retrying to parse with ; delimiter", err)
		// try to re-parse with ; delimiter
		file.Seek(0, 0)

		reader = csv.NewReader(file)
		reader.Comma = ';'

		records, err = reader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}
	}

	hasHeader := false
	if options.HasHeader {
		hasHeader = true
	} else if options.HasNoHeader {
		hasHeader = false
	} else {
		// first row might be a header, guess
		hasEmpty := false
		for _, col := range records[0] {
			if col == "" {
				hasEmpty = true
				break
			}
		}

		if !hasEmpty {
			log.Println("First row has no empty column, assuming it is the header.")
			hasHeader = true
		}
	}

	log.Println(records[0])

	var bank *cfg.Bank
	var exists bool

	if options.BankName != "" {
		bank, exists = config.Banks[options.BankName]
		if !exists {
			log.Fatalf("Bank with name %s not found in config", options.BankName)
		}
	} else {
		// try to get bank config by filename match
		baseFile := filepath.Base(fileName)
		for _, b := range config.Banks {
			if b.FileNamePattern != "" {
				match, _ := regexp.MatchString(b.FileNamePattern, baseFile)
				if match {
					bank = b
					log.Printf("Bank determined by file name pattern `%s`", b.FileNamePattern)
					break
				}
			}
		}

		if !hasHeader {
			log.Fatal("CVS file does not contain header row and bank name was not provided.  Cannot determine bank configuration.")
		}

		// determine bank automatically
		if bank == nil {
			bank, exists = cfg.GetBankConfig(records[0], config.Banks)
			if !exists {
				log.Fatalf("No configured bank matches the cvs file %s", fileName)
			}
		}

		log.Printf("Using automatically detected bank %s", bank.Name)
	}

	if (bank.ColumnIndices == cfg.ColumnIndices{}) {
		if !hasHeader {
			log.Fatal("ColumnIndices not present in bank config and there is not header row to determine them from names.")
		}

		bank.ColumnIndices = bank.NamesToIndices(records[0])
	}

	var transactions []t.Transaction
	for i, record := range records {
		if i == 0 && hasHeader {
			continue
		}

		trans := t.FromCsvRecord(record, config, bank)
		transactions = append(transactions, trans)
	}

	return transactions, bank
}

func main() {
	var options Options
	var parser = flags.NewParser(&options, flags.Default)

	args, err := parser.Parse()

	if err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	config := cfg.LoadConfig(options.Config)
	config.ValidateConfig()

	transactions, bank := readCsv(args[0], options, config)
	bank.ValidateBankConfig()

	buffer := t.TransactionBuffer{Transactions: []t.Transaction{}}
	unknownPayees := make([]string, 1)

	for _, trans := range transactions {
		if trans.IsIgnored() {
			continue
		}

		payee, exists := trans.GetPayee()
		if !exists {
			if !Contains(unknownPayees, payee.Name) {
				unknownPayees = append(unknownPayees, payee.Name)
			}
		}

		if twinType := trans.IsTwinTransaction(); twinType != nil {
			buffer.Twin = *twinType
			buffer.Append(trans)
		} else {
			if bank := trans.IsTransactionToOwnAccount(); bank != nil {
				// only generate outgoing payments between our own accounts
				if trans.AmountAccount > 0 {
					continue
				}
			}

			fmt.Println(trans.FormatTrans(buffer))
			buffer = t.TransactionBuffer{Transactions: []t.Transaction{}}
		}
	}

	fmt.Fprintf(os.Stderr, "\n\n%s", s.Join(unknownPayees, "\n"))
}
