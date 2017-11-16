package transactor

// Transactor
// Bench
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	// "sync/atomic"
	"testing"
)

func BenchmarkCreditSingle(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 9223372036854775).End()
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tr.Begin().Credit(int64(uint16(i)), "USD", 1).End()
	}
}

func BenchmarkCreditParallel(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	tr.AddUnit(1234567)
	tr.Begin().Debit(1234567, "USD", 9223372036854775806).End()
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 9223372036854775).End()
	}

	i := uint16(0)
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tr.Begin().Credit(int64(i), "USD", 1).End()
			i++
		}
	})
}

func BenchmarkDebitSingle(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(int64(uint16(i)), "USD", 1).End()
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tr.Begin().Debit(int64(uint16(i)), "USD", 1).End()
	}
}

func BenchmarkDebitParallel(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	tr.AddUnit(1234567)
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(int64(uint16(i)), "USD", 1).End()
	}

	i := uint16(0)
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tr.Begin().Debit(int64(i), "USD", 1).End()
			i++
		}
	})
}

func BenchmarkTransferSingle(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 100000).End()
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tr.Begin().Credit(int64(uint16(i)), "USD", 1).Debit(int64(uint16(i+1)), "USD", 1).End()
	}
}

func BenchmarkTransferParallel(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	tr.AddUnit(1234567)
	tr.AddUnit(1234568)
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 100000).End()
	}

	i := uint16(0)
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tr.Begin().Credit(int64(i), "USD", 1).Debit(int64(i+1), "USD", 1).End()
			i++
		}
	})
}

func BenchmarkBuySingle(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 100000).End()
		tr.Begin().Debit(i, "APPLE", 5000).End()
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tr.Begin().
			Credit(int64(uint16(i)), "USD", 10).Debit(int64(uint16(i+1)), "USD", 10).
			Debit(int64(uint16(i)), "APPLE", 2).Credit(int64(uint16(i+1)), "APPLE", 2).
			End()
	}
}

func BenchmarkBuyParallel(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	tr.AddUnit(1234567)
	tr.AddUnit(1234568)
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 100000).End()
		tr.Begin().Debit(i, "APPLE", 5000).End()
	}

	i := uint16(0)
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tr.Begin().
				Credit(int64(uint16(i)), "USD", 10).Debit(int64(uint16(i+1)), "USD", 10).
				Debit(int64(uint16(i)), "APPLE", 2).Credit(int64(uint16(i+1)), "APPLE", 2).
				End()
			i++
		}
	})
}

/*

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
*/
