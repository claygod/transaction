package transactor

// Transactor
// Account
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"errors"
	//"log"
	"runtime"
	//"sync"
	"sync/atomic"
)

type Account struct {
	counter int64
	balance int64
}

// newAccount - create new account.
func newAccount(amount int64) *Account {
	k := &Account{balance: amount}
	return k
}

func (a *Account) addition(amount int64) int64 {
	b := atomic.LoadInt64(&a.balance)
	nb := b + amount

	if nb < 0 || atomic.CompareAndSwapInt64(&a.balance, b, nb) {
		return nb
	}

	for i := trialLimit; i > trialStop; i-- {
		b := atomic.LoadInt64(&a.balance)
		nb := b + amount
		if nb < 0 || atomic.CompareAndSwapInt64(&a.balance, b, nb) {
			return nb
		}
		//if nb < 0 {
		//	return nb
		//}
		//if atomic.CompareAndSwapInt64(&a.balance, b, nb) {
		//	return nb
		//}
		runtime.Gosched()
	}
	return permitError
}

func (a *Account) total() int64 {
	return atomic.LoadInt64(&a.balance)
}

// catch - ловим разрешение на проведение операций с аккаунтом
func (a *Account) catch() bool {
	if atomic.LoadInt64(&a.counter) < 0 {
		return false
	}
	if atomic.AddInt64(&a.counter, 1) > 0 {
		return true
	}
	atomic.AddInt64(&a.counter, -1)
	return false
}

func (a *Account) throw() {
	atomic.AddInt64(&a.counter, -1)
}

func (a *Account) start() bool {
	var currentCounter int64
	for i := trialLimit; i > trialStop; i-- {
		currentCounter = atomic.LoadInt64(&a.counter)
		if currentCounter >= 0 {
			return true
		}
		// the variable `currentCounter` is expected to be `permitError`
		if atomic.CompareAndSwapInt64(&a.counter, permitError, 0) {
			return true
		}
		runtime.Gosched()
	}
	return false
}

func (a *Account) stop() bool {
	var currentCounter int64

	for i := trialLimit; i > trialStop; i-- {
		currentCounter = atomic.LoadInt64(&a.counter)
		switch {
		case currentCounter == 0:
			if atomic.CompareAndSwapInt64(&a.counter, 0, permitError) {
				return true
			}
		case currentCounter > 0:
			atomic.CompareAndSwapInt64(&a.counter, currentCounter, currentCounter+permitError)
		case currentCounter == permitError:
			return true
		}
		runtime.Gosched()
	}
	currentCounter = atomic.LoadInt64(&a.counter)
	if currentCounter < 0 && currentCounter > permitError {
		atomic.AddInt64(&a.counter, -permitError)
	}
	return false
}

func (a *Account) stopUnsafe() {
	atomic.StoreInt64(&a.counter, permitError)
}
