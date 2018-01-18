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

func TestTransactorStart(t *testing.T) {
	//trialLimit = 200
	tr := New()

	if !tr.Start() {
		t.Error("Now the start is possible!")
	}
	//tr.Stop()
	//t.Error(tr.Stop())
	//t.Error(tr.hasp)
	//t.Error(tr.counter)
	//t.Error(stateClosed)

	//trialLimit = trialLimitConst
}

func TestTransactorStop(t *testing.T) {
	trialLimit = 200
	tr := New()

	if !tr.Start() {
		t.Error("Now the start is possible!")
	}

	if !tr.Stop() {
		t.Error("Now the stop is possible!")
	}

	trialLimit = trialLimitConst
}

/*
 */
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

	if tr.AddUnit(123) == Ok {
		t.Error("You can not re-add a unit")
	}
}

func TestTransactorDelUnit(t *testing.T) {
	tr := New()
	tr.Start()
	if _, err := tr.DelUnit(123); err == Ok {
		t.Error("Removed non-existent unit")
	}

	tr.AddUnit(123)

	if _, err := tr.DelUnit(123); err != Ok {
		t.Error("The unit has not been deleted")
	}
}

func TestTransactorTotallUnit(t *testing.T) {
	tr := New()
	tr.Start()
	tr.AddUnit(123)
	tr.Begin().Debit(123, "USD", 1).End()

	arr, err := tr.TotalUnit(123)

	if err != Ok {
		t.Error("Failed to get information on the unit")
	}

	if b, ok := arr["USD"]; ok != true || b != 1 {
		t.Error("The received information on the unit is erroneous")
	}

	if _, err := tr.TotalUnit(456); err == Ok {
		t.Error("A unit does not exist, there must be an error")
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
	tr.Stop()
}

func TestTransactorLoad(t *testing.T) {
	path := "./test.tdb"
	tr := New()
	tr.Start()
	tr.AddUnit(123)
	tr.Begin().Debit(123, "USD", 7).End()
	// --- tr.AddUnit(456)
	// -- tr.Begin().Debit(456, "USD", 12).End()
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
