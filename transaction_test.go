package transactor

// Transactor
// Test
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

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

	tn := tr.Begin()
	tn = tn.Debit(123, "USD", 5)
	if res := tn.exeTransaction(); res != Ok {
		t.Error(res)
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
