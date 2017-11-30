package transactor

// Transactor
// Account test
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"testing"
)

func TestAccountAdd(t *testing.T) {
	a := newAccount(100)
	if a.debit(50) != 150 {
		t.Error("Error adding")
	}
}

func TestAccountCredit(t *testing.T) {
	a := newAccount(100)
	if a.credit(50) != 50 {
		t.Error("Account incorrectly performed a credit operation")
	}
	if a.credit(51) >= 0 {
		t.Error("There must be a negative result")
	}
}

func TestAccountDebit(t *testing.T) {
	a := newAccount(100)
	if a.debit(50) != 150 {
		t.Error("Account incorrectly performed a debit operation")
	}
}

func TestAccountTotal(t *testing.T) {
	a := newAccount(100)
	a.credit(1)
	a.debit(1)
	if a.total() != 100 {
		t.Error("Balance error")
	}
}

func TestAccountCath(t *testing.T) {
	a := newAccount(100)
	if !a.catch() {
		t.Error("Account is free, but it was not possible to catch it")
	}
	if a.catch(); a.counter != 2 {
		t.Error("Account counter error")
	}
}

func TestAccountThrow(t *testing.T) {
	a := newAccount(100)
	a.catch()
	if a.throw(); a.counter != 0 {
		t.Error("Failed to decrement the counter")
	}
}

func TestAccountStart(t *testing.T) {
	a := newAccount(100)
	a.counter = permitError
	if !a.start() {
		t.Error("Could not launch this account")
	}
}

func TestAccountStop(t *testing.T) {
	a := newAccount(100)
	if !a.stop() {
		t.Error("Could not stop this account")
	}
	// This part of the test should be run only with a small value of the constant "trialLimit"
	// a.counter = 1
	// if a.stop() {
	// 	t.Error("Could not stop this account22")
	// }
}
