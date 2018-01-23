package transactor

// Transactor
// Bench
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	// "sync/atomic"
	"testing"
)

func BenchmarkCreditSequence(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 223372036854775).End()
	}
	u := 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tr.Begin().Credit(int64(uint16(i)), "USD", 1).End()

		//reqs := []*Request{
		//	&Request{id: int64(uint16(u)), key: "USD", amount: -1},
		//}
		//newTransaction2(&tr, reqs).exeTransaction()
		u++
	}
}

func BenchmarkCreditParallel(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	tr.AddUnit(1234567)
	tr.Begin().Debit(1234567, "USD", 223372036854775806).End()
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 223372036854775).End()
	}

	u := 0
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			//tr.Begin().Credit(int64(uint16(i)), "USD", 1).End()
			//i++
			//reqs := []*Request{
			//	&Request{id: int64(uint16(u)), key: "USD", amount: -1},
			//}
			//newTransaction2(&tr, reqs).exeTransaction()
			tr.Begin().Credit(int64(uint16(u)), "USD", 1).End()
			u++
		}
	})
}

func BenchmarkDebitSequence(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 1).End()
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

	i := 0
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tr.Begin().Debit(int64(uint16(i)), "USD", 1).End()
			i++
		}
	})
}

func BenchmarkTransferSequence(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 100000000).End()
	}
	u := 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tr.Begin().Credit(int64(uint16(u)), "USD", 1).Debit(int64(uint16(u+1)), "USD", 1).End()
		u += 2
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
		tr.Begin().Debit(i, "USD", 100000000).End()
	}

	u := 0
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tr.Begin().Credit(int64(uint16(u)), "USD", 1).Debit(int64(uint16(u+1)), "USD", 1).End()
			u += 2
		}
	})
}

/*
func BenchmarkBuyUnsafeSequence(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 100000000).End()
		tr.Begin().Debit(i, "APPLE", 5000000).End()
	}

	u := 0
	//var tn *Transaction

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		reqs := []*Request{
			&Request{id: int64(uint16(u)), key: "USD", amount: -10},
			&Request{id: int64(uint16(u + 1)), key: "APPLE", amount: -2},
			&Request{id: int64(uint16(u + 1)), key: "USD", amount: 10},
			&Request{id: int64(uint16(u)), key: "APPLE", amount: 2},
		}
		//Unsafe(&tr, reqs)
		tr.Unsafe(reqs)
		u += 2
	}
}

func BenchmarkBuyUnsafenParallel(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	tr.AddUnit(1234567)
	tr.AddUnit(1234568)
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 100000000).End()
		tr.Begin().Debit(i, "APPLE", 5000000).End()
	}

	u := 0
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			reqs := []*Request{
				&Request{id: int64(uint16(u)), key: "USD", amount: -10},
				&Request{id: int64(uint16(u + 1)), key: "APPLE", amount: -2},
				&Request{id: int64(uint16(u + 1)), key: "USD", amount: 10},
				&Request{id: int64(uint16(u)), key: "APPLE", amount: 2},
			}
			//Unsafe(&tr, reqs)
			tr.Unsafe(reqs)
			u += 2
		}
	})
}

*/
func BenchmarkBuySequence(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 100000000).End()
		tr.Begin().Debit(i, "APPLE", 5000000).End()
	}

	//tn := tr.Begin().
	//	Credit(int64(uint16(u)), "USD", 10).Debit(int64(uint16(u+1)), "USD", 10).
	//	Debit(int64(uint16(u)), "APPLE", 2).Credit(int64(uint16(u+1)), "APPLE", 2)

	b.StartTimer()
	u := 0
	for i := 0; i < b.N; i++ {
		//tn.reqs[0].account
		tr.Begin().
			Credit(int64(uint16(u)), "USD", 10).Debit(int64(uint16(u+1)), "USD", 10).
			Debit(int64(uint16(u)), "APPLE", 2).Credit(int64(uint16(u+1)), "APPLE", 2).
			End()
		u += 2
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
		tr.Begin().Debit(i, "USD", 100000000).End()
		tr.Begin().Debit(i, "APPLE", 5000000).End()
	}

	u := 0
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tr.Begin().
				Credit(int64(uint16(u)), "USD", 10).Debit(int64(uint16(u+1)), "USD", 10).
				Debit(int64(uint16(u)), "APPLE", 2).Credit(int64(uint16(u+1)), "APPLE", 2).
				End()
			u += 2
		}
	})
}

/*
func BenchmarkBuyPreparation2Parallel(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	tr.AddUnit(1234567)
	tr.AddUnit(1234568)
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 100000000).End()
		tr.Begin().Debit(i, "APPLE", 5000000).End()
	}

	u := 0
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			reqs := []*Request{
				&Request{id: int64(uint16(u)), key: "USD", amount: -10},
				&Request{id: int64(uint16(u + 1)), key: "USD", amount: 10},
				&Request{id: int64(uint16(u)), key: "APPLE", amount: 2},
				&Request{id: int64(uint16(u + 1)), key: "APPLE", amount: -2},
			}
			newTransaction2(&tr, reqs)

			u += 2
		}
	})
}

func BenchmarkBuyPreparation1Parallel(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	tr.AddUnit(1234567)
	tr.AddUnit(1234568)
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 100000000).End()
		tr.Begin().Debit(i, "APPLE", 5000000).End()
	}

	u := 0
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			newTransaction(&tr).
				//tr.Begin().
				Credit(int64(uint16(u)), "USD", 10).Debit(int64(uint16(u+1)), "USD", 10).
				Debit(int64(uint16(u)), "APPLE", 2).Credit(int64(uint16(u+1)), "APPLE", 2)
			u += 2
		}
	})
}
*/
func BenchmarkTrGetAccount2Sequence(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 100000000).End()
		tr.Begin().Debit(i, "APPLE", 5000000).End()
	}
	u := 0
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		//tn.reqs[0].account
		tr.getAccount(int64(uint16(u)), "USD")
		u += 2
	}
}

func BenchmarkTrGetAccount2Parallel(b *testing.B) {
	b.StopTimer()

	tr := New()
	tr.Start()
	for i := int64(0); i < 65536; i++ {
		tr.AddUnit(i)
		tr.Begin().Debit(i, "USD", 100000000).End()
		tr.Begin().Debit(i, "APPLE", 5000000).End()
	}
	u := 0

	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tr.getAccount(int64(uint16(u)), "USD")
			u += 2
		}
	})
}

func BenchmarkUnitGetAccountSequence(b *testing.B) {
	b.StopTimer()

	un := newUnit()

	tr := New()
	tr.Start()
	for i := int64(0); i < 65536; i++ {
		un.getAccount("USD")
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		//tn.reqs[0].account
		un.getAccount("USD")
	}
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
