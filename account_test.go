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

	if a.credit(1) < 0 {
		t.Error("A positive result was expected")
	}
	//if a.credit(1) == permitError {
	//	t.Error("Account must not be blocked! The number of waiting cycles exceeded.")
	//}
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
	a.counter = -1
	if a.catch() {
		t.Error("There must be an answer `false`")
	}
	a.counter = 1
	if !a.catch() {
		t.Error("There must be an answer `true`")
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
	trialLimit = 200
	a.counter = -1
	if a.start() {
		t.Error("Could not launch this account")
	}
	a.counter = 1
	if !a.start() {
		t.Error("The account already has a positive counter and the function should return the `true`")
	}
	a.counter = 0
	if !a.start() {
		t.Error("The account is already running and the function should return the` true`")
	}
	a.counter = permitError
	if !a.start() {
		t.Error("Account counter is in a position that allows launch. The launch is not carried out erroneously.")
	}
	trialLimit = trialLimitConst
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
