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
	hasp    int64
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

/*
func (a *Account) getBalance() int64 {
	return atomic.LoadInt64(&a.balance)
}

func (a *Account) getDebt() int64 {
	return atomic.LoadInt64(&a.debt)
}
*/
func (a *Account) permit() bool {
	for i := trialLimit; i > 5; i-- {
		if atomic.LoadInt64(&a.hasp) < 0 {
			return false
		}
		if atomic.CompareAndSwapInt64(&a.hasp, 0, 1) {
			return true
		}
		runtime.Gosched()
	}
	return false
}

type AccountState [2]int64
