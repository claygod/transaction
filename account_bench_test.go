package transactor

// Transactor
// Account bench
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	// "sync/atomic"
	"testing"
)

// +++++++

func BenchmarkAccountCreditAtomFreeOk(b *testing.B) {
	b.StopTimer()
	a := newAccount(9223372036854775807)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.creditAtomicFree(1)
	}
}

func BenchmarkAccountCreditAtomFreeErr(b *testing.B) {
	b.StopTimer()
	a := newAccount(1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.creditAtomicFree(1)
	}
}

func BenchmarkAccountDebitAtomFreeOk(b *testing.B) {
	b.StopTimer()
	a := newAccount(1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.debitAtomicFree(1)
	}
}

// -==== --

func BenchmarkAccountTotal(b *testing.B) {
	b.StopTimer()
	a := newAccount(4767567567)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.total()
	}
}

// --
