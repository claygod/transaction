package transaction

// Core
// Test
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

func TestUsage(t *testing.T) {
	path := "./test.tdb"
	tr := New()
	if !tr.Start() {
		t.Error("Now the start is possible!")
	}

	tr.AddUnit(123)
	tr.Begin().Debit(123, "USD", 7777).End()
	tr.Save(path)

	if _, err := tr.DelUnit(123); err == Ok {
		t.Error(err)
	}

	if err := tr.Begin().Debit(123, "USD", 7777).End(); err != Ok {
		t.Error(err)
	}

	tr.Stop()
	os.Remove(path)
}

func TestTransfer(t *testing.T) {
	tr := New()
	tr.Start()
	tr.AddUnit(14760464)
	tr.AddUnit(2674560)

	if err := tr.Begin().Debit(14760464, "USD", 11).End(); err != Ok {
		t.Error(err)
	}
	if err := tr.Begin().Debit(2674560, "USD", 7).End(); err != Ok {
		t.Error(err)
	}
	if err := tr.Begin().Credit(2674560, "USD", 2).End(); err != Ok {
		t.Error(err)
	}

	if err := tr.Begin().
		Credit(2674560, "USD", 4).
		Debit(14760464, "USD", 4).
		End(); err != Ok {
		t.Error(err)
	}
}

func TestCoreStart(t *testing.T) {
	tr := New()

	if !tr.Start() {
		t.Error("Now the start is possible!")
	}
	tr.Stop()
	trialLimit = trialStop
	if tr.Start() {
		t.Error("Now the start is possible!")
	}

	trialLimit = trialLimitConst
}

func TestCoreStop(t *testing.T) {
	trialLimit = 200
	tr := New()

	if !tr.Start() {
		t.Error("Now the start is possible!")
	}

	if ok, _ := tr.Stop(); !ok {
		t.Error("Now the stop is possible!")
	}
	tr.Start()
	trialLimit = trialStop
	//tr.hasp = stateClosed
	if ok, _ := tr.Stop(); ok {
		t.Error("Due to the limitation of the number of iterations, stopping is impossible")
	}

	trialLimit = trialLimitConst
}

func TestCoreGetAccount(t *testing.T) {
	tr := New()

	if _, err := tr.getAccount(123, "USD"); err == Ok {
		t.Error("We must get an error!")
	}

	tr.AddUnit(123)

	if _, err := tr.getAccount(123, "USD"); err != Ok {
		t.Error("We should not get an error!")
	}
}

func TestCoreAddUnit(t *testing.T) {
	tr := New()
	tr.Start()

	if tr.AddUnit(123) != Ok {
		t.Error("Unable to add unit")
	}

	if tr.AddUnit(123) == Ok {
		t.Error("You can not re-add a unit")
	}

	tr.hasp = stateClosed

	if tr.AddUnit(456) == Ok {
		t.Error("Due to the blocking, it was not possible to add a new unit.")
	}
}

func TestCoreDelUnit(t *testing.T) {
	tr := New()
	tr.Start()
	if _, err := tr.DelUnit(123); err == Ok {
		t.Error("Removed non-existent unit")
	}

	tr.AddUnit(123)

	if _, err := tr.DelUnit(123); err != Ok {
		t.Error("The unit has not been deleted")
	}

	tr.AddUnit(456)
	tr.Begin().Debit(456, "USD", 5).End()

	tr.storage.getUnit(456).getAccount("USD").counter = 1

	trialLimit = trialStop + 100
	if _, err := tr.DelUnit(456); err == Ok {
		t.Error("The unit has not been deleted")
	}

	tr.hasp = stateClosed

	if _, err := tr.DelUnit(456); err == Ok {
		t.Error("Due to the blocking, it was not possible to del a unit.")
	}
	trialLimit = trialLimitConst
}

func TestCoreTotalUnit(t *testing.T) {
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

	tr.hasp = stateClosed
	if _, err := tr.TotalUnit(123); err != ErrCodeCoreCatch {
		t.Error("Resource is locked and can not allow operation")
	}
	tr.hasp = stateOpen
}

func TestCoreTotalAccount(t *testing.T) {
	tr := New()
	tr.Start()
	tr.AddUnit(123)
	tr.Begin().Debit(123, "USD", 1).End()

	balance, err := tr.TotalAccount(123, "USD")

	if err != Ok {
		t.Error("Failed to get information on the account")
	}

	if balance != 1 {
		t.Error("The received information on the account is erroneous")
	}

	if _, err := tr.TotalAccount(456, "USD"); err == Ok {
		t.Error("A account does not exist, there must be an error")
	}

	if balance, err := tr.TotalAccount(123, "EUR"); err != Ok || balance != 0 {
		t.Error("A account does not exist, there must be an error")
	}

	tr.hasp = stateClosed
	if _, err := tr.TotalAccount(123, "USD"); err != ErrCodeCoreCatch {
		t.Error("Resource is locked and can not allow operation")
	}
	tr.hasp = stateOpen
}

func TestCoreSave(t *testing.T) {
	path := "./test.tdb"
	tr := New()
	tr.Start()
	tr.AddUnit(123)
	tr.Begin().Debit(123, "USD", 7).End()

	trialLimit = trialStop
	tr.hasp = stateClosed
	if tr.Save(path) == Ok {
		t.Error("The lock should prevent the file from being saved")
	}
	trialLimit = trialLimitConst

	tr.hasp = stateOpen
	if tr.Save(path) != Ok {
		t.Error("There is no lock, saving should be successful")
	}

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

func TestCoreLoad(t *testing.T) {
	path := "./test.tdb"
	pathFake := "./testFake.tdb"
	tr := New()
	tr.Start()
	tr.AddUnit(123)
	tr.Begin().Debit(123, "USD", 7).End()
	tr.Save(path)
	tr.Stop()
	tr2 := New()
	if res, _ := tr2.Load(pathFake); res == Ok {
		t.Errorf("The file `%s` does not exist", pathFake)
	}
	if res, _ := tr2.Load(path); res != Ok {
		t.Errorf("The file `%s` does exist", pathFake)
	}
	tr2.Start()
	if res, _ := tr2.Load(path); res != Ok {
		t.Errorf("Error loading the database file (%d)", res)
	}
	balance, res := tr2.TotalAccount(123, "USD")
	if balance != 7 {
		t.Errorf("Error in account balance (%d)", balance)
	}
	if res != Ok {
		t.Errorf("Error in the downloaded account (%d)", res)
	}
	os.Remove(path)
}
