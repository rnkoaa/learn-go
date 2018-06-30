package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

func main() {
	jsonStr := `{"data": {
    "object": "card",
    "id": "card_123",
    "last4": "4242"
  }}`
	var m map[string]Data
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		panic(err)
	}
	fmt.Println(m)
	data := m["data"]
	if data.Card != nil {
		fmt.Println(data.Card)
	}
	if data.BankAccount != nil {
		fmt.Println(data.BankAccount)
	}

	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

type BankAccount struct {
	ID            string `json:"id"`
	Object        string `json:"object"`
	RoutingNumber string `json:"routing_number"`
}

type Card struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Last4  string `json:"last4"`
}

type Data struct {
	*Card
	*BankAccount
}

func (d Data) MarshalJSON() ([]byte, error) {
	if d.Card != nil {
		return json.Marshal(d.Card)
	} else if d.BankAccount != nil {
		return json.Marshal(d.BankAccount)
	} else {
		return json.Marshal(nil)
	}
}

func (d *Data) UnmarshalJSON(data []byte) error {
	temp := struct {
		Object string `json:"object"`
	}{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	if temp.Object == "card" {
		var c Card
		if err := json.Unmarshal(data, &c); err != nil {
			return err
		}
		d.Card = &c
		d.BankAccount = nil
	} else if temp.Object == "bank_account" {
		var ba BankAccount
		if err := json.Unmarshal(data, &ba); err != nil {
			return err
		}
		d.BankAccount = &ba
		d.Card = nil
	} else {
		return errors.New("Invalid object value")
	}
	return nil
}
