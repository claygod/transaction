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
	//a.counter = 1
	//if a.stop() {
	//	t.Error("Could not stop this account22")
	//}
}

/*
func TestAccountReserveOk(t *testing.T) {
	a := newAccount(100)
	if a.reserve(50) != nil {
		t.Error("The available funds for the reservation were sufficient")
	}
	if a.debt != 50 {
		t.Error("The expected value is 50, and the value obtained is ", a.debt)
	}
}

func TestAccountReserveError(t *testing.T) {
	a := newAccount(100)
	if a.reserve(150) == nil {
		t.Error("The available funds for the reservation were not enough")
	}
	if a.debt != 0 || a.balance != 100 {
		t.Error("The expected value is 0, and the value obtained is ", a.debt)
	}
}

func TestAccountUnreserve(t *testing.T) {
	a := newAccount(100)
	a.reserve(50)
	if a.unreserve(10) != nil {
		t.Error("Error unblocking blocked funds")
	}
	if a.unreserve(50) == nil {
		t.Error("Unblocking funds beyond the limit")
	}
}

func TestAccountUnreserveTotal(t *testing.T) {
	a := newAccount(100)
	a.reserve(50)
	a.reserve(50)
	if a.unreserveTotal() != nil {
		t.Error("Unable to unlock all funds")
	}
	if a.debt != 0 || a.balance != 100 {
		t.Error("Unblocked funds are not transferred to the balance")
	}

}

func TestAccountGive(t *testing.T) {
	a := newAccount(150)
	a.reserve(50)
	if a.give(50) != nil {
		t.Error("The blocked amount was not enough")
	}
	if a.give(50) == nil {
		t.Error("Paid with non-existent funds")
	}
}

func TestAccountDelOk(t *testing.T) {
	a := newAccount(100)
	a.topup(50)
	if a.withdraw(150) != nil {
		t.Error("The funds available on the balance sheet were sufficient")
	}
}

func TestAccountDelError(t *testing.T) {
	a := newAccount(100)
	a.topup(50)
	if a.withdraw(200) == nil {
		t.Error("The funds available on the balance sheet were insufficient")
	}
}
*/
