package transaction

// Transaction
// Bench
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"testing"
)

func BenchmarkSpeedMutex(b *testing.B) {
	b.StopTimer()
	n := newNode()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		n.m.Lock()
		n.m.Unlock()
	}
}

func BenchmarkSpeedAtom(b *testing.B) {
	b.StopTimer()
	n := newNode()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		n.lock()
		//n.hasp = 0
		n.unlock()
	}
}

func BenchmarkFreezeUnfreezeNodeadlock(b *testing.B) {
	b.StopTimer()
	k := New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		k.TransactionStart(uint64(i), uint64(i)+1)
		k.TransactionEnd(uint64(i), uint64(i)+1)
	}
}

func BenchmarkFreezeUnfreezeDeadlock(b *testing.B) {
	b.StopTimer()
	iterat := 1000
	k := New()
	for i := 0; i < iterat; i++ {
		go k.TransactionStart(uint64(i), uint64(i)+1)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		k.TransactionEnd(uint64(i), uint64(i)+1)
	}
}

func BenchmarkMapAdd(b *testing.B) {
	b.StopTimer()
	m := make(map[uint64]bool)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m[uint64(uint8(i))] = true
	}
}
