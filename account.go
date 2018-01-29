package transactor

// Core
// Account
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"runtime"
	"sync/atomic"
)

/*
Account - keeps a balance.
Account balance can not be less than zero.

The counter has two tasks:
	counts the number of operations performed
	stops the account (new operations are not started)
*/
type Account struct {
	counter int64
	balance int64
}

/*
newAccount - create new account.
*/
func newAccount(amount int64) *Account {
	k := &Account{balance: amount}
	return k
}

/*
addition - add to the balance the input variable.
If the result of adding the balance and the input
variable is less than zero, the balance does not change.

Returned result:
	greater than or equal to zero // OK
	less than zero // Error
*/
func (a *Account) addition(amount int64) int64 {

	// The hidden part of the code allows you
	// to speed up a bit by avoiding the cycle start
	//b := atomic.LoadInt64(&a.balance)
	//nb := b + amount
	//if nb < 0 || atomic.CompareAndSwapInt64(&a.balance, b, nb) {
	//	return nb
	//}

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

/*
total - current account balance.
*/
func (a *Account) total() int64 {
	return atomic.LoadInt64(&a.balance)
}

/*
catch - get permission to perform operations with the account.
*/
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

/*
throw - current account operation has been completed
*/
func (a *Account) throw() {
	atomic.AddInt64(&a.counter, -1)
}

/*
start - start an account.
*/
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

/*
stop - stop an account.
*/
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
