package main

import (
	"container/list"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/shopspring/decimal"
)

type transaction struct {
	recipient           string
	transaction_type    string
	transaction_details string
	amount              decimal.Decimal
}

type rule struct {
	recipient   string
	description string
	category    string
}

func open_file(file_name string) *os.File {
	// open file
	f, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func load_rules(file_name string) *list.List {
	f := open_file(file_name)

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	rules := list.New()

	//Load rules
	for {
		rec, err := csvReader.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err, f.Name())
		}

		rules.PushBack(rule{recipient: rec[0], description: rec[1], category: rec[2]})
	}

	f.Close()
	return rules
}

func load_transactions(file_name string) *list.List {
	f := open_file(file_name)

	csvReader := csv.NewReader(f)
	transactions := list.New()
	//Load transactions
	for {
		rec, err := csvReader.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// do something with read line
		amount, err := decimal.NewFromString(rec[5])
		if err != nil {
			//log.Default(err)
		}
		transactions.PushBack(transaction{recipient: rec[1], transaction_type: rec[3], transaction_details: rec[4], amount: amount})
	}

	f.Close()
	return transactions
}

func get_category(rules *list.List, recipient string, details string) string {
	for e := rules.Front(); e != nil; e = e.Next() {
		rule := e.Value.(rule)

		if strings.Contains(recipient, rule.recipient) && (rule.description == "" || strings.Contains(details, rule.description)) {
			return rule.category
		}
	}

	fmt.Printf("%+v, %+v\n", recipient, details)

	return "Uncategorized"
}

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Print("Usage: ./main [transaction_csv] [rules_csv]\n")
		os.Exit(3)
	}

	rules := load_rules(args[1])
	fmt.Printf("Loaded %+v rules\n", rules.Len())

	transactions := load_transactions(args[0])
	fmt.Printf("Loaded %+v transactions\n", transactions.Len())

	fmt.Print("--------------------------\n")

	summary := map[string]decimal.Decimal{}

	for e := transactions.Front(); e != nil; e = e.Next() {
		transaction := e.Value.(transaction)
		cat := get_category(rules, transaction.recipient, transaction.transaction_details)

		res, exists := summary[cat]

		if exists {
			summary[cat] = res.Add(transaction.amount)
		} else {
			summary[cat] = transaction.amount
		}
	}

	sum, err := decimal.NewFromString("0.0")

	if err != nil {
		//log.Default(err)
	}

	for i, e := range summary {
		fmt.Printf("[%+v] %+v EUR\n", i, e)

		sum = sum.Add(e)
	}

	fmt.Print("--------------------------\n")
	fmt.Printf("[TOTAL] %+v EUR\n", sum)
}
