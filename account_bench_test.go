package transactor

// Transactor
// Account bench
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	// "sync/atomic"
	"testing"
)

// +++++++
/*
func BenchmarkChannelIn(b *testing.B) {
	b.StopTimer()
	ch := make(chan int, 100000005)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ch <- 1
	}
}

func BenchmarkChannelOut(b *testing.B) {
	b.StopTimer()
	ch := make(chan int, 100000005)
	for i := 0; i < 100000000; i++ {
		ch <- 1
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		<-ch
	}
}


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
*/
// -==== --

func BenchmarkAccountTotal(b *testing.B) {
	b.StopTimer()
	a := newAccount(4767567567)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.total()
	}
}

func BenchmarkAccountStartStop(b *testing.B) {
	b.StopTimer()
	a := newAccount(4767567567)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.start()
		a.stop()
	}
}

func BenchmarkAccountAddition(b *testing.B) {
	b.StopTimer()
	a := newAccount(1)
	a.start()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.addition(1)
	}
}

// --
