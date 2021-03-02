package transaction

// Core
// Account bench
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"testing"
)

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
