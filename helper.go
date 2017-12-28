package transactor

// Transactor
// Helper
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"errors"
	//"log"
	"runtime"
	//"sync"
	"sync/atomic"
)

/*
func credit(balance *int64, amount int64) int64 {
	for i := trialLimit; i > trialStop; i-- {
		b := atomic.LoadInt64(balance)
		nb := b - amount
		// if nb < 0 || atomic.CompareAndSwapInt64(balance, b, nb) {
		//	return nb
		// }
		if nb < 0 {
			return nb
		}
		if atomic.CompareAndSwapInt64(balance, b, nb) {
			return nb
		}
		runtime.Gosched()
	}
	return permitError
}

func debit(balance *int64, amount int64) int64 {
	return atomic.AddInt64(balance, amount)
}

func total(balance *int64) int64 {
	return atomic.LoadInt64(balance)
}
*/
func catch(counter *int64) bool {
	if atomic.LoadInt64(counter) < 0 {
		return false
	}
	if atomic.AddInt64(counter, 1) > 0 {
		return true
	}
	atomic.AddInt64(counter, -1)
	return false
}

func throw(counter *int64) {
	atomic.AddInt64(counter, -1)
}

func start(counter *int64) bool {
	var currentCounter int64
	for i := trialLimit; i > trialStop; i-- {
		currentCounter = atomic.LoadInt64(counter)
		if currentCounter >= 0 {
			return true
		}
		// the variable `currentCounter` is expected to be `permitError`
		if atomic.CompareAndSwapInt64(counter, permitError, 0) {
			return true
		}
		runtime.Gosched()
	}
	return false
}

func stop(counter *int64) bool {
	var currentCounter int64

	for i := trialLimit; i > trialStop; i-- {
		currentCounter = atomic.LoadInt64(counter)
		switch {
		case currentCounter == 0:
			if atomic.CompareAndSwapInt64(counter, 0, permitError) {
				return true
			}
		case currentCounter > 0:
			atomic.CompareAndSwapInt64(counter, currentCounter, currentCounter+permitError)
		case currentCounter == permitError:
			return true
		}
		runtime.Gosched()
	}
	currentCounter = atomic.LoadInt64(counter)
	if currentCounter < 0 && currentCounter > permitError {
		atomic.AddInt64(counter, -permitError)
	}
	return false
}

func stopUnsafe(counter *int64) {
	atomic.StoreInt64(counter, permitError)
}

// type AccountState [2]int64
