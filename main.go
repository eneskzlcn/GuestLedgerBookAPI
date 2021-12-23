package main

import (
	"errors"
	"strings"
)

type GuestLedger struct{
	Email string `json:"email"`
	Message string `json:"message"`
}

var GuestLedgerBook []GuestLedger

func findGuestLedgerByEmail(email string) (interface{},error) {
	for _,value := range GuestLedgerBook{
		if strings.Compare(value.Email,email) == 0{
			return value, nil
		}
	}
	return nil,errors.New("guest ledger of given email not found")
}