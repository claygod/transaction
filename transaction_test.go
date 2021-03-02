package transaction

// Core
// Test
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
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

	tr.hasp = stateClosed

	if tr.Begin().Debit(123, "USD", 5).End() != ErrCodeCoreCatch {
		t.Error("Resource is locked and can not allow operation")
	}

	tr.hasp = stateOpen

	if tr.Begin().Debit(123, "USD", 5).End() != Ok {
		t.Error("Error executing a transaction")
	}

	tr.storage.getUnit(123).getAccount("USD").counter = -1

	if tr.Begin().Debit(123, "USD", 5).End() != ErrCodeTransactionCatch {
		t.Error("The transaction could not cath the account")
	}

	tr.storage.delUnit(123)

	if tr.Begin().Debit(123, "USD", 5).End() == Ok {
		t.Error("The requested unit does not exist")
	}
}

func TestTransactionCatch(t *testing.T) {
	tr := New()
	tr.Start()
	tr.AddUnit(1)
	tr.AddUnit(2)
	tr.Begin().Debit(1, "USD", 5).End()
	tr.Begin().Debit(2, "USD", 5).End()

	tn := tr.Begin().Credit(1, "USD", 1)
	tn.reqs[0].account = tr.storage.getUnit(1).getAccount("USD")

	if tn.catch() != Ok {
		t.Error("TestTransactionCatch")
	}

	tr.storage.getUnit(1).getAccount("USD").counter = -1

	tn2 := tr.Begin().Credit(1, "USD", 2)
	tn2.reqs[0].account = tr.storage.getUnit(1).getAccount("USD")

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
