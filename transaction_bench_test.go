package transactor

// Transactor
// Bench
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"sync/atomic"
	"testing"
)

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
