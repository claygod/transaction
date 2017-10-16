package transaction

// Transaction
// Account
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"errors"
	"runtime"
	"sync/atomic"
)

type account struct {
	balance int64
	debt    int64
}

// newAccount - create new account.
func newAccount(amount int64) account {
	k := account{balance: amount}
	return k
}

func (a *account) add(amount int64) int64 {
	b := atomic.AddInt64(&a.balance, amount)
	return b
}

func (a *account) reserve(amount int64) error {
	if atomic.AddInt64(&a.balance, -(amount)) < 0 {
		atomic.AddInt64(&a.balance, amount)
		return errors.New("Insufficient funds in the account")
	}
	atomic.AddInt64(&a.debt, amount)
	return nil
}

func (a *account) unreserve(amount int64) error {
	if atomic.AddInt64(&a.debt, -(amount)) < 0 {
		atomic.AddInt64(&a.debt, amount)
		return errors.New("So much was not reserved")
	}
	atomic.AddInt64(&a.balance, amount)
	return nil
}

func (a *account) unreserveTotal() error {
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

func (a *account) give(amount int64) error {
	if atomic.AddInt64(&a.debt, -(amount)) < 0 {
		atomic.AddInt64(&a.debt, amount)
		return errors.New("So many not reserved")
	}
	return nil
}

func (a *account) del(amount int64) error {
	if atomic.AddInt64(&a.balance, -(amount)) < 0 {
		atomic.AddInt64(&a.balance, amount)
		return errors.New("So many not reserved")
	}
	return nil
}

func (a *account) getBalance() int64 {
	return atomic.LoadInt64(&a.balance)
}

func (a *account) getDebt() int64 {
	return atomic.LoadInt64(&a.debt)
}
