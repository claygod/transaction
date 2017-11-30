package transactor

// Transactor
// Account test
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"testing"
)

func TestUnitGetAccount(t *testing.T) {
	u := newUnit()
	if u.getAccount("USD") == nil {
		t.Error("New account not received")
	}
}

func TestUnitTotal(t *testing.T) {
	u := newUnit()
	if len(u.total()) != 0 {
		t.Error("The unit has phantom accounts")
	}
	u.getAccount("USD").debit(5)
	all := u.total()
	if num, ok := all["USD"]; !ok || num != 5 {
		t.Error("Lost data from one of the accounts")
	}
}

func TestUnitDelAccount(t *testing.T) {
	u := newUnit()
	if u.delAccount("USD") != ErrCodeAccountNotExist {
		t.Error("Deleted non-existent account")
	}
	u.getAccount("USD").debit(5)
	if u.delAccount("USD") != ErrCodeAccountNotEmpty {
		t.Error("Deleted account with non-zero balance")
	}
	u.getAccount("USD").credit(5)
	if u.delAccount("USD") != Ok {
		t.Error("Account not deleted, although it is possible")
	}
}

func TestUnitDelAllAccount(t *testing.T) {
	u := newUnit()
	if lst, err := u.delAllAccounts(); len(lst) != 0 || err != Ok {
		t.Error("Error deleting all accounts (and they are not)")
	}
	// This part of the test should be run only with a small value of the constant "trialLimit"
	// u.getAccount("USD").counter = 1
	// if lst, err := u.delAllAccounts(); len(lst) != 1 || err != ErrCodeAccountNotStop {
	// 	t.Error("When cleaning a unit, it removes a non-empty account")
	// }
	u.getAccount("USD").counter = 0
	u.getAccount("USD").debit(5)
	if lst, err := u.delAllAccounts(); len(lst) != 1 || err != ErrCodeUnitNotEmpty {
		t.Error("When cleaning a unit, it removes a non-empty account")
	}
}

func TestUnitDel(t *testing.T) {
	u := newUnit()
	u.getAccount("USD")
	if lst := u.del(); len(lst) != 0 {
		t.Error("It was not possible to delete all accounts (but this is possible)")
	}
	u.getAccount("USD").debit(5)
	if lst := u.del(); len(lst) == 0 {
		t.Error("It turned out to delete all accounts (but this is impossible)")
	}
}

func TestUnitStop(t *testing.T) {
	u := newUnit()
	u.getAccount("USD").debit(5)
	if lst := u.stop(); len(lst) != 0 {
		t.Error("I could stop all accounts (but it's impossible)")
	}
	// This part of the test should be run only with a small value of the constant "trialLimit"
	// u.getAccount("USD").counter = 1
	// if lst := u.stop(); len(lst) == 0 {
	// 	t.Error("I could not stop all accounts (but it's possible)")
	// }
}

func TestUnitStart(t *testing.T) {
	u := newUnit()
	u.getAccount("USD").debit(5)
	u.stop()
	if lst := u.start(); len(lst) != 0 {
		t.Error("I could start all accounts (but it's impossible)")
	}
}
