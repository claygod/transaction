package transactor

// Transactor
// Account
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"errors"
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

func (a *Account) state() (int64, int64) {
	if !a.permit() {
		return -1, -1
	}
	b1 := a.balance
	d1 := a.debt
	atomic.StoreInt64(&a.hasp, 0)

	return b1, d1
}

func (a *Account) creditAtomic2(amount int64) int64 {
	for i := trialLimit; i > trialStop; i-- {
		if atomic.CompareAndSwapInt64(&a.hasp, 0, 1) {
			b := a.balance - amount
			if b >= 0 {
				a.balance = b
			}
			atomic.StoreInt64(&a.hasp, 0)
			return b
		}
		runtime.Gosched()
	}
	return permitError

}

func (a *Account) creditAtomic(amount int64) int64 {
	if !a.permit() {
		return permitError
	}
	b := a.balance - amount
	if b >= 0 {
		a.balance = b
	}
	atomic.StoreInt64(&a.hasp, 0)
	return b
}

func (a *Account) creditAtomicFree(amount int64) int64 {
	for i := trialLimit; i > trialStop; i-- {
		b := atomic.LoadInt64(&a.balance)
		nb := b - amount
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

func (a *Account) credit(amount int64) int64 {
	a.Lock()
	b := a.balance - amount
	if b >= 0 {
		a.balance = b
	}
	a.Unlock()
	return b
}

func (a *Account) debit(amount int64) int64 {
	a.Lock()
	b := a.balance + amount
	a.balance = b
	a.Unlock()
	return b
}

func (a *Account) debitAtomic2(amount int64) int64 {
	for i := trialLimit; i > trialStop; i-- {
		if atomic.CompareAndSwapInt64(&a.hasp, 0, 1) {
			b := a.balance + amount
			a.balance = b
			atomic.StoreInt64(&a.hasp, 0)
			return b
		}
		runtime.Gosched()
	}
	return permitError
}

func (a *Account) debitAtomic(amount int64) int64 {
	if !a.permit() {
		return permitError
	}
	b := a.balance + amount
	a.balance = b
	atomic.StoreInt64(&a.hasp, 0)
	return b
}

func (a *Account) total() int64 {
	return atomic.LoadInt64(&a.balance)
}

func (a *Account) topup(amount int64) int64 {
	if !a.permit() {
		return -1
	}
	a.balance += amount
	b := a.balance
	atomic.StoreInt64(&a.hasp, 0)
	return b
}

func (a *Account) reserve(amount int64) error {
	if !a.permit() {
		return errors.New("Account not available")
	}
	if a.balance-amount < 0 {
		atomic.StoreInt64(&a.hasp, 0)
		return errors.New("Insufficient funds in the account")
	}
	a.balance -= amount
	a.debt += amount
	atomic.StoreInt64(&a.hasp, 0)
	return nil
}

func (a *Account) unreserve(amount int64) error {
	if !a.permit() {
		return errors.New("Account not available")
	}
	if a.debt-amount < 0 {
		atomic.StoreInt64(&a.hasp, 0)
		return errors.New("Insufficient funds in the account")
	}
	a.balance += amount
	a.debt -= amount
	atomic.StoreInt64(&a.hasp, 0)
	return nil
}

func (a *Account) unreserveTotal() error {
	if !a.permit() {
		return errors.New("Account not available")
	}
	a.balance += a.debt
	a.debt = 0
	atomic.StoreInt64(&a.hasp, 0)
	return nil
}

func (a *Account) give(amount int64) error {
	if !a.permit() {
		return errors.New("Account not available")
	}
	if a.debt-amount < 0 {
		atomic.StoreInt64(&a.hasp, 0)
		return errors.New("So many not reserved")
	}
	a.debt -= amount
	atomic.StoreInt64(&a.hasp, 0)
	return nil
}

func (a *Account) withdraw(amount int64) error {
	if !a.permit() {
		return errors.New("Account not available")
	}
	if a.balance-amount < 0 {
		atomic.StoreInt64(&a.hasp, 0)
		return errors.New("So many not reserved")
	}
	a.balance -= amount
	atomic.StoreInt64(&a.hasp, 0)
	return nil
}

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
	atomic.AddInt64(&a.counter, -1)
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
