package transactor

// Transactor
// Test
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"fmt"
	"testing"
)

func TestTransfer(t *testing.T) {
	tr := New()
	//tr.Load("test.tdb")
	tr.Start()
	//ta.Transfer().From(1).To(2).Account("ABC").Count(5).Do()
	tr.AddUnit(14760464)
	tr.AddUnit(2674560)
	/*
		for i := int64(1); i < 10; i++ {
			tr.AddUnit(i)
			if err := tr.Begin().Debit(i, "USD", 2).End(); err != ErrOk {
				t.Error(err)
			}
		}
	*/
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
	tr.Save("test.tdb")
	//tr.Load("test.tdb")

	if err := tr.Begin().
		Credit(2674560, "USD", 4).
		Debit(14760464, "USD", 4).
		End(); err != Ok {
		t.Error(err)
	}

}