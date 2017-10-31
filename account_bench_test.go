package transaction

// Transaction
// Account bench
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	// "sync/atomic"
	"testing"
)

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
