package main

type GuestLedger struct{
	Email string `json:"email"`
	Message string `json:"message"`
}

var GuestLedgerBook []GuestLedger

