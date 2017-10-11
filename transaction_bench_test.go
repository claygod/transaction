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
