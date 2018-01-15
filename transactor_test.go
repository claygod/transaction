package transactor

// Transactor
// Test
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

func TestTransfer(t *testing.T) {
	tr := New()
	//tr.Load("test.tdb")
	tr.Start()
	//ta.Transfer().From(1).To(2).Account("ABC").Count(5).Do()
	tr.AddUnit(14760464)
	tr.AddUnit(2674560)

	if err := tr.Begin().Debit(14760464, "USD", 11).End(); err != Ok {
		t.Error(err)
	}
	if err := tr.Begin().Debit(2674560, "USD", 7).End(); err != Ok {
		t.Error(err)
	}
	//t.Error("------", ta.Begin().Credit(2674560, "USD", 9).End())
	if err := tr.Begin().Credit(2674560, "USD", 2).End(); err != Ok {
		t.Error(err)
	}
	//tr.Save("test.tdb")
	//tr.Load("test.tdb")

	if err := tr.Begin().
		//Op(2674560, "USD", -4).
		//Op(14760464, "USD", 4).
		Credit(2674560, "USD", 4).
		Debit(14760464, "USD", 4).
		End(); err != Ok {
		t.Error(err)
	}

}

func TestTransactorUnsafe(t *testing.T) {
	tr := New()
	tr.Start()
	tr.AddUnit(123)

	reqs := []*Request{
		&Request{id: 123, key: "USD", amount: 10},
	}
	//Unsafe(&tr, reqs)

	if tr.Unsafe(reqs) != Ok {
		t.Error("Transaction error in Unsafe mode")
	}

	if res, _ := tr.TotalAccount(123, "USD"); res != 10 {
		t.Error("Invalid transaction result")
	}
}

func TestTransactorAddUnit(t *testing.T) {
	tr := New()
	tr.Start()
	tr.AddUnit(123)

	reqs := []*Request{
		&Request{id: 123, key: "USD", amount: 10},
	}
	//Unsafe(&tr, reqs)

	if tr.Unsafe(reqs) != Ok {
		t.Error("Transaction error in Unsafe mode")
	}

	if res, _ := tr.TotalAccount(123, "USD"); res != 10 {
		t.Error("Invalid transaction result")
	}
}

func TestTransactorSave(t *testing.T) {
	path := "./test.tdb"
	tr := New()
	tr.Start()
	tr.AddUnit(123)
	tr.Begin().Debit(123, "USD", 7).End()
	tr.Save(path)

	endLine := []byte(endLineSymbol)
	separator := []byte(separatorSymbol)

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Can not find saved file")
	}
	str := bytes.Split(bs, endLine)[0]
	a := bytes.Split(str, separator)
	if len(a) != 3 {
		t.Error("Invalid number of columns")
	}
	id, err := strconv.ParseInt(string(a[0]), 10, 64)
	if err != nil {
		t.Error("Error converting string to integer (id account)")
	}
	balance, err := strconv.ParseInt(string(a[1]), 10, 64)
	if err != nil {
		t.Error("Error converting string to integer (balance account)")
	}
	if id != 123 {
		t.Error("The account identifier does not match")
	}
	if balance != 7 {
		t.Error("The account balance does not match")
	}
	os.Remove(path)
}

func TestTransactorLoad(t *testing.T) {
	path := "./test.tdb"
	tr := New()
	tr.Start()
	tr.AddUnit(123)
	tr.Begin().Debit(123, "USD", 7).End()
	//tr.AddUnit(456)
	//tr.Begin().Debit(456, "USD", 12).End()
	tr.Save(path)
	tr.Stop()

	tr2 := New()
	tr2.Load(path)
	tr2.Start()
	if res := tr2.Load(path); res != Ok {
		t.Error(fmt.Sprintf("Error loading the database file (%d)", res))
	}
	balance, res := tr2.TotalAccount(123, "USD")
	if balance != 7 {
		t.Error(fmt.Sprintf("Error in account balance (%d)", balance))
	}
	if res != Ok {
		t.Error(fmt.Sprintf("Error in the downloaded account (%d)", res))
	}

	// os.Remove(path)
}
