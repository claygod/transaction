package transaction

// Transaction
// Account
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"errors"
	"runtime"
	"sync/atomic"
)

type Account struct {
	balance int64
	debt    int64
}

// newAccount - create new account.
func newAccount(amount int64) *Account {
	k := &Account{balance: amount}
	return k
}

func (a *Account) state() (int64, int64) {
	var b1, d1, b2, d2 int64
	b1 = atomic.LoadInt64(&a.balance)
	d1 = atomic.LoadInt64(&a.debt)

	for {
		b2 = atomic.LoadInt64(&a.balance)
		d2 = atomic.LoadInt64(&a.debt)
		if b1 == b2 && d1 == d2 {
			break
		}
		b1 = b2
		d1 = d2
		runtime.Gosched()
	}
	return b1, d1
}

func (a *Account) add(amount int64) int64 {
	b := atomic.AddInt64(&a.balance, amount)
	return b
}

func (a *Account) reserve(amount int64) error {
	if atomic.AddInt64(&a.balance, -(amount)) < 0 {
		atomic.AddInt64(&a.balance, amount)
		return errors.New("Insufficient funds in the account")
	}
	atomic.AddInt64(&a.debt, amount)
	return nil
}

func (a *Account) unreserve(amount int64) error {
	if atomic.AddInt64(&a.debt, -(amount)) < 0 {
		atomic.AddInt64(&a.debt, amount)
		return errors.New("So much was not reserved")
	}
	atomic.AddInt64(&a.balance, amount)
	return nil
}

func (a *Account) unreserveTotal() error {
	var d int64
	for i := trialLimit; i > 0; i-- {
		d = atomic.LoadInt64(&a.debt)
		if atomic.CompareAndSwapInt64(&a.debt, d, 0) == true {
			atomic.AddInt64(&a.balance, d)
			return nil
		}
		runtime.Gosched()
	}
	return errors.New("Unable to unblock funds")
}

func (a *Account) give(amount int64) error {
	if atomic.AddInt64(&a.debt, -(amount)) < 0 {
		atomic.AddInt64(&a.debt, amount)
		return errors.New("So many not reserved")
	}
	return nil
}

func (a *Account) del(amount int64) error {
	if atomic.AddInt64(&a.balance, -(amount)) < 0 {
		atomic.AddInt64(&a.balance, amount)
		return errors.New("So many not reserved")
	}
	return nil
}

func (a *Account) getBalance() int64 {
	return atomic.LoadInt64(&a.balance)
}

func (a *Account) getDebt() int64 {
	return atomic.LoadInt64(&a.debt)
}
