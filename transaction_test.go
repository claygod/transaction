package transaction

// Transaction
// Test
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	//"runtime"
	//"sync"
	"sync/atomic"
	"testing"
)

func TestUnfreezeUnfrozen(t *testing.T) {
	k := New()
	k.TransactionStart(101, 102)
	if k.TransactionEnd(101, 200) == nil {
		t.Error("The program unlocked a number that was not frozen")
	}
}

/*
func TestDoNotUnlock(t *testing.T) {
	k := New()
	k.TransactionStart(101, 102)
	if k.TransactionStart(101, 200) != false {
		t.Error("Double lock of the same number")
	}
}
*/
/*
func TestDoNotUnlockParallel(t *testing.T) {
	k := New()
	go k.TransactionStart(101, 102)
	go k.TransactionStart(101, 200)

	go k.TransactionEnd(101, 200)
	go k.TransactionEnd(101, 102)

	if k.TransactionEnd(101, 200) == nil {
		t.Error("Parallel threads prevented the program from terminating correctly.")
	}
}
*/
func TestFreeze(t *testing.T) {
	//var wg sync.WaitGroup
	iterat := 90000
	k := New()
	for i := 0; i < iterat; i++ {

		k.TransactionStart(uint64(i), uint64(i)+1)

		go k.TransactionEnd(uint64(i), uint64(i)+1)
		//wg.Add(1)
	}
	//for i := 0; i < 10; i++ {
	//	wg.Add(1)
	//	go k.TransactionEnd(uint64(i), uint64(i)+1)
	//}

	for i := 0; i < iterat; i++ {
		//runtime.Gosched()
	}

	//wg.Wait()
	if count := atomic.LoadInt64(&k.counter); count != 0 {
		t.Error(fmt.Sprintf("Received `%d` instead of `0`", count))
	}
}
