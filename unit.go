package transaction

// Core
// Unit
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"errors"
	//"log"
	"sync"
)

/*
unit - aggregates accounts.
*/
type unit struct {
	sync.Mutex
	accounts map[string]*Account
}

/*
newUnit - create new Unit.
*/
func newUnit() *unit {
	k := &unit{accounts: make(map[string]*Account)}
	return k
}

/*
getAccount - take account by key.
If there is no such account, it will be created (with zero balance).
*/
func (u *unit) getAccount(key string) *Account {
	a, ok := u.accounts[key]
	if !ok {
		u.Lock()
		defer u.Unlock()

		a, ok = u.accounts[key]
		if !ok {
			a = newAccount(0)
			u.accounts[key] = a
		}
	}
	return a
}

/*
total - current balance of all accounts.
At the time of the formation of the answer will be a lock.
*/
func (u *unit) total() map[string]int64 {
	t := make(map[string]int64)
	//u.Lock()
	for k, a := range u.accounts {
		t[k] = a.total()
	}
	//u.Unlock()
	return t
}

/*
totalUnsafe - current balance of all accounts (fast and unsafe).
Without locking. Use this method only in serial (non-parallel) mode.
*/
func (u *unit) totalUnsafe() map[string]int64 {
	t := make(map[string]int64)
	for k, a := range u.accounts {
		t[k] = a.total()
	}
	return t
}

/*
delAccount - delete account.
The account can not be deleted if it is not stopped or has a non-zero balance.
*/
func (u *unit) delAccount(key string) errorCodes {
	u.Lock()
	defer u.Unlock()
	a, ok := u.accounts[key]
	if !ok {
		return ErrCodeAccountNotExist
	}
	if a.total() != 0 {
		return ErrCodeAccountNotEmpty
	}
	if !a.stop() {
		return ErrCodeAccountNotStop
	}

	delete(u.accounts, key)

	return Ok
}

/*
func (u *Unit) delAccountUnsafe(key string) errorCodes {
	u.Lock()
	defer u.Unlock()
	_, ok := u.accounts[key]
	if !ok {
		return ErrCodeAccountNotExist
	}
	delete(u.accounts, key)

	return Ok
}
*/

/*
delAllAccounts - delete all accounts.
On error, a list of non-stopped or non-empty accounts is returned.

Returned codes:
	ErrCodeAccountNotStop // not stopped
	ErrCodeUnitNotEmpty // not empty
	Ok
*/
func (u *unit) delAllAccounts() ([]string, errorCodes) {
	u.Lock()
	defer u.Unlock()
	if notStop := u.stop(); len(notStop) != 0 {
		return notStop, ErrCodeAccountNotStop
	}
	if notDel := u.delStoppedAccounts(); len(notDel) != 0 {
		u.start() // Undeleted accounts are restarted
		return notDel, ErrCodeUnitNotEmpty
	}

	return nil, Ok
}

/*
delStoppedAccounts - delete all accounts (they are stopped).
Returns a list of not deleted accounts (with a non-zero balance).
*/
func (u *unit) delStoppedAccounts() []string {
	notDel := make([]string, 0, len(u.accounts))
	for k, a := range u.accounts {
		if a.total() != 0 {
			notDel = append(notDel, k)
		} else {
			delete(u.accounts, k)
		}
	}
	return notDel
}

/*
start - start all accounts.
Returns a list of not starting accounts.
*/
func (u *unit) start() []string {
	notStart := make([]string, 0, len(u.accounts))
	for k, a := range u.accounts {
		if !a.start() {
			notStart = append(notStart, k)
		}
	}
	return notStart
}

/*
stop - stop all accounts.
Returns a list of non-stopped accounts.
*/
func (u *unit) stop() []string {
	notStop := make([]string, 0, len(u.accounts))
	for k, a := range u.accounts {
		if !a.stop() {
			notStop = append(notStop, k)
		}
	}
	return notStop
}
