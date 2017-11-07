package transactor

// Transactor
// Account bench
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	// "sync/atomic"
	"testing"
)

func BenchmarkAccountCreditOk(b *testing.B) {
	b.StopTimer()
	a := newAccount(9223372036854775807)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.credit(1)
	}
}

func BenchmarkAccountCreditErr(b *testing.B) {
	b.StopTimer()
	a := newAccount(1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.credit(1)
	}
}

func BenchmarkAccountCreditAtomOk(b *testing.B) {
	b.StopTimer()
	a := newAccount(9223372036854775807)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.creditAtomic(1)
	}
}

func BenchmarkAccountCreditAtomErr(b *testing.B) {
	b.StopTimer()
	a := newAccount(1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.creditAtomic(1)
	}
}

func BenchmarkAccountCreditAtom2Ok(b *testing.B) {
	b.StopTimer()
	a := newAccount(9223372036854775807)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.creditAtomic2(1)
	}
}

func BenchmarkAccountCreditAtom2Err(b *testing.B) {
	b.StopTimer()
	a := newAccount(1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.creditAtomic2(1)
	}
}

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

func BenchmarkAccountDebitOk(b *testing.B) {
	b.StopTimer()
	a := newAccount(1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.debit(1)
	}
}

func BenchmarkAccountDebitAtomOk(b *testing.B) {
	b.StopTimer()
	a := newAccount(1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.debitAtomic(1)
	}
}

func BenchmarkAccountDebitAtom2Ok(b *testing.B) {
	b.StopTimer()
	a := newAccount(1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.debitAtomic2(1)
	}
}

// ---

func BenchmarkAccountAdd(b *testing.B) {
	b.StopTimer()
	a := newAccount(100)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.topup(1)
	}
}

func BenchmarkAccountAddParallel(b *testing.B) {
	b.StopTimer()
	a := newAccount(100)
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			a.topup(1)
		}
	})
}

func BenchmarkAccountReserve(b *testing.B) {
	b.StopTimer()
	a := newAccount(9223372036854775807)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.reserve(1)
	}
}

func BenchmarkAccountReserveParallel(b *testing.B) {
	b.StopTimer()
	a := newAccount(9223372036854775807)
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			a.reserve(1)
		}
	})
}

func BenchmarkAccountGive(b *testing.B) {
	b.StopTimer()
	a := newAccount(9223372036854775807)
	a.reserve(9223372036854775807)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.give(1)
	}
}

func BenchmarkAccountGiveParallel(b *testing.B) {
	b.StopTimer()
	a := newAccount(9223372036854775807)
	a.reserve(9223372036854775807)
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			a.give(1)
		}
	})
}
