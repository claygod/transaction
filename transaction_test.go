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
	//ta.Transfer().From(1).To(2).Account("ABC").Count(5).Do()
	tr.AddCustomer(14760464)
	tr.AddCustomer(2674560)

	if err := tr.Begin().Debit(2674560, "USD", 7).End(); err != nil {
		t.Error(err)
	}
	//t.Error("------", ta.Begin().Credit(2674560, "USD", 9).End())
	if err := tr.Begin().Credit(2674560, "USD", 3).End(); err != nil {
		t.Error(err)
	}

	if err := tr.Begin().
		Credit(2674560, "USD", 4).
		Debit(14760464, "USD", 4).
		End(); err != nil {
		t.Error(err)
	}

}
