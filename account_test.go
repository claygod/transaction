package transaction

// Core
// Account test
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"testing"
)

func TestAccountAdd(t *testing.T) {
	a := newAccount(100)

	if a.addition(50) != 150 {
		t.Error("Error addition (The sum does not match)")
	}

	if a.addition(-200) != -50 {
		t.Error("Error adding (Negative balance)")
	}

	trialLimit = trialStop
	if a.addition(-200) != permitError {
		t.Error("Error adding")
	}

	trialLimit = trialLimitConst
}

func TestAccountAddingOverflow1(t *testing.T) {
	a := newAccount(-1)

	if res := a.addition(-0x8000000000000000); res != -0x8000000000000000 {
		t.Error("Overflow error (1)")
		t.Error(res)
	}
}

func TestAccountAddingOverflow2(t *testing.T) {
	a := newAccount(-0x8000000000000000)

	if res := a.addition(-5); res != -5 {
		t.Error("Overflow error (2)")
		t.Error(res)
	}
}

func TestAccountCredit(t *testing.T) {
	a := newAccount(100)

	if a.addition(-50) != 50 {
		t.Error("Account incorrectly performed a credit operation")
	}

	if a.addition(-51) >= 0 {
		t.Error("There must be a negative result")
	}

	if a.addition(-1) < 0 {
		t.Error("A positive result was expected")
	}

	if b := a.addition(-1); b != 48 {
		t.Error("Awaiting 48 and received:", b)
	}
}

func TestAccountDebit(t *testing.T) {
	a := newAccount(100)

	if a.addition(50) != 150 {
		t.Error("Account incorrectly performed a debit operation")
	}
}

func TestAccountTotal(t *testing.T) {
	a := newAccount(100)
	a.addition(-1)
	a.addition(1)

	if a.total() != 100 {
		t.Error("Balance error", a.total())
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
	trialLimit = 200

	if !a.stop() {
		t.Error("Could not stop this account")
	}

	a.counter = 1

	if a.stop() {
		t.Error("Could not stop this account22")
	}

	trialLimit = trialLimitConst
}
