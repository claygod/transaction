package transactor

// Core
// Account test
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

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
	u.getAccount("USD").addition(5)
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
	u.getAccount("USD").addition(5)
	if u.delAccount("USD") != ErrCodeAccountNotEmpty {
		t.Error("Deleted account with non-zero balance")
	}
	u.getAccount("USD").addition(-5)
	if u.delAccount("USD") != Ok {
		t.Error("Account not deleted, although it is possible")
	}
}

func TestUnitDelAllAccount(t *testing.T) {
	u := newUnit()
	if lst, err := u.delAllAccounts(); len(lst) != 0 || err != Ok {
		t.Error("Error deleting all accounts (and they are not)")
	}
	trialLimit = trialStop + 100
	u.getAccount("USD").counter = 1
	if lst, err := u.delAllAccounts(); len(lst) != 1 || err != ErrCodeAccountNotStop {
		t.Error("When cleaning a unit, it removes a non-empty account")
	}
	trialLimit = trialLimitConst

	u.getAccount("USD").counter = 0
	u.getAccount("USD").addition(5)
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
	u.getAccount("USD").addition(5)
	if lst := u.del(); len(lst) == 0 {
		t.Error("It turned out to delete all accounts (but this is impossible)")
	}
}

func TestUnitStop(t *testing.T) {
	u := newUnit()
	u.getAccount("USD").addition(5)
	if lst := u.stop(); len(lst) != 0 {
		t.Error("I could stop all accounts (but it's impossible)")
	}
	trialLimit = trialStop + 100
	u.getAccount("USD").counter = 1
	if lst := u.stop(); len(lst) == 0 {
		t.Error("I could not stop all accounts (but it's possible)")
	}
	trialLimit = trialLimitConst
}

func TestUnitStart(t *testing.T) {
	u := newUnit()
	u.getAccount("USD").addition(5)
	u.stop()
	if lst := u.start(); len(lst) != 0 {
		t.Error("I could start all accounts (but it's impossible)")
	}

	//u.getAccount("USD").counter = 1
	//trialLimit = trialStop + 100
	//if lst := u.start(); len(lst) != 1 {
	//	t.Error(len(lst))
	//}
	//trialLimit = trialLimitConst
}
