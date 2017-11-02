package transaction

// Transaction
// Bench
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"sync/atomic"
	"testing"
)

/*
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
*/

func BenchmarkMapRead(b *testing.B) {
	b.StopTimer()
	m := make(map[int]bool)
	for i := 0; i < 70000; i++ {
		m[i] = true
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = m[int(int16(i))]
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

func BenchmarkSliceAdd(b *testing.B) {
	b.StopTimer()
	m := make([]bool, 0, 7000000)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m = append(m, true)
	}
}

func BenchmarkCAS(b *testing.B) {
	b.StopTimer()
	var m int64 = 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		atomic.CompareAndSwapInt64(&m, 0, 0)
	}
}

func BenchmarkAtomicStore(b *testing.B) {
	b.StopTimer()
	var m int64 = 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		atomic.StoreInt64(&m, 0)
	}
}

func BenchmarkAtomicLoad(b *testing.B) {
	b.StopTimer()
	var m int64 = 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		atomic.LoadInt64(&m)
	}
}

func BenchmarkAtomicAdd(b *testing.B) {
	b.StopTimer()
	var m int64 = 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		atomic.AddInt64(&m, 1)
	}
}
