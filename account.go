package transactor

// Transactor
// Account
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"errors"
	//"log"
	"runtime"
	"sync"
	"sync/atomic"
)

type Account struct {
	sync.Mutex
	hasp    int64
	counter int64
	balance int64
	debt    int64
}

// newAccount - create new account.
func newAccount(amount int64) *Account {
	k := &Account{balance: amount}
	return k
}

func (a *Account) creditAtomicFree(amount int64) int64 {
	for i := trialLimit; i > trialStop; i-- {
		b := atomic.LoadInt64(&a.balance)
		nb := b - amount
		//if nb > b { // variable overflow
		//	return permitError
		//}
		if nb < 0 || atomic.CompareAndSwapInt64(&a.balance, b, nb) {
			return nb
		}
		runtime.Gosched()
	}
	return permitError
}

func (a *Account) debitAtomicFree(amount int64) int64 {
	return atomic.AddInt64(&a.balance, amount)
}

func (a *Account) total() int64 {
	return atomic.LoadInt64(&a.balance)
}

// --

func (a *Account) permit() bool {
	for i := trialLimit; i > trialStop; i-- {
		if atomic.CompareAndSwapInt64(&a.hasp, 0, 1) {
			return true
		}
		runtime.Gosched()
	}
	return false
}

func (a *Account) catch() bool {
	if atomic.LoadInt64(&a.counter) < 0 {
		return false
	}
	if atomic.AddInt64(&a.counter, 1) > 0 {
		return true
	}
	atomic.AddInt64(&a.counter, permitError/2)
	return false
}

func (a *Account) throw() {
	atomic.AddInt64(&a.counter, -1)
}

func (a *Account) start() bool {
	for i := trialLimit; i > trialStop; i-- {
		c := atomic.LoadInt64(&a.counter)
		if c >= 0 {
			return true
		}
		if atomic.CompareAndSwapInt64(&a.counter, c, 0) {
			return true
		}
		runtime.Gosched()
	}
	return false
}

func (a *Account) stop() bool {
	for i := trialLimit; i > trialStop; i-- {
		if atomic.LoadInt64(&a.counter) < 0 || atomic.CompareAndSwapInt64(&a.counter, 0, permitError) {
			return true
		}
		runtime.Gosched()
	}
	return false
}

type AccountState [2]int64
