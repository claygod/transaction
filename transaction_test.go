package transaction

// Core
// Test
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"fmt"
	"testing"
)

func TestCreditPrepare(t *testing.T) {
	tr := New()
	tr.Start()
	tn := tr.Begin()
	tn = tn.Credit(123, "USD", 5)
	if len(tn.reqs) != 1 {
		t.Error("When preparing a transaction, the credit operation is lost.")
	}
	if tn.reqs[0].amount != -5 {
		t.Error("In the lending operation, the amount.")
	}
	if tn.reqs[0].id != 123 {
		t.Error("In the lending operation, the ID.")
	}
	if tn.reqs[0].key != "USD" {
		t.Error("In the lending operation, the name of the value.")
	}
}

func TestTransactionExe(t *testing.T) {
	tr := New()
	tr.Start()
	tr.AddUnit(123)

	//tn := tr.Begin().Debit(123, "USD", 5)

	if tr.Begin().Debit(123, "USD", 5).End() != Ok {
		t.Error("Error executing a transaction")
	}

	tr.storage.getUnit(123).accounts["USD"].counter = -1

	if tr.Begin().Debit(123, "USD", 5).End() != ErrCodeTransactionCatch {
		t.Error("The transaction could not cath the account")
	}

	tr.storage.delUnit(123)
	if tr.Begin().Debit(123, "USD", 5).End() == Ok {
		t.Error("The requested unit does not exist")
	}

	//tr.storage.getUnit(123).accounts["USD"].counter = 0
	//tr.counter = stateClosed

}

func TestTransactionCatch(t *testing.T) {
	tr := New()
	tr.Start()
	tr.AddUnit(1)
	tr.AddUnit(2)
	tr.Begin().Debit(1, "USD", 5).End()
	tr.Begin().Debit(2, "USD", 5).End()

	tn := tr.Begin().Credit(1, "USD", 1)
	tn.reqs[0].account = tr.storage.getUnit(1).accounts["USD"]
	if tn.catch() != Ok {
		t.Error("TestTransactionCatch")
	}
	tr.storage.getUnit(1).accounts["USD"].counter = -1

	tn2 := tr.Begin().Credit(1, "USD", 2)
	tn2.reqs[0].account = tr.storage.getUnit(1).accounts["USD"]

	if tn2.catch() == Ok {
		t.Error("TestTransactionCatch 222")
	}
}

func TestTransactionRollback(t *testing.T) {
	tr := New()
	tr.Start()
	tr.AddUnit(123)
	tr.Begin().Debit(123, "USD", 7).End()

	tn := newTransaction(&tr)
	tn.Credit(123, "USD", 1)
	tn.fill()
	//tn.catch()
	tn.exeTransaction()
	if num, _ := tr.TotalAccount(123, "USD"); num != 6 {
		t.Error("Credit operation is not carried out")
	}
	tn.rollback(1)
	if num, _ := tr.TotalAccount(123, "USD"); num != 7 {
		t.Error("Not rolled back")
	}

}

/*
func TestDebitPrepare(t *testing.T) {
	tr := New()
	tr.Start()
	tn := tr.Begin()
	tn = tn.Debit(123, "USD", 5)
	if len(tn.up) != 1 {
		t.Error("When preparing a transaction, the debit operation is lost.")
	}
	if tn.up[0].amount != 5 {
		t.Error("In the debit operation, the amount.")
	}
	if tn.up[0].id != 123 {
		t.Error("In the debit operation, the ID.")
	}
	if tn.up[0].key != "USD" {
		t.Error("In the debit operation, the name of the value.")
	}
}
*/
