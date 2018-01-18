package transactor

// Transactor
// Unit
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"errors"
	//"log"
	"sync"
)

type Unit struct {
	sync.Mutex
	accounts map[string]*Account
}

// newUnit - create new Unit.
func newUnit() *Unit {
	k := &Unit{accounts: make(map[string]*Account)}
	return k
}

func (u *Unit) getAccount(key string) *Account {
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

func (u *Unit) total() map[string]int64 {
	t := make(map[string]int64)
	u.Lock()
	for k, a := range u.accounts {
		t[k] = a.total()
	}
	u.Unlock()
	return t
}

func (u *Unit) totalUnsave() map[string]int64 {
	t := make(map[string]int64)
	for k, a := range u.accounts {
		t[k] = a.total()
	}
	return t
}

func (u *Unit) delAccount(key string) errorCodes {
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

func (u *Unit) delAllAccounts() ([]string, errorCodes) {
	u.Lock()
	defer u.Unlock()
	if notStop := u.stop(); len(notStop) != 0 {
		return notStop, ErrCodeAccountNotStop
	}
	if notDel := u.del(); len(notDel) != 0 {
		return notDel, ErrCodeUnitNotEmpty
	}

	return nil, Ok
}

func (u *Unit) del() []string {
	notDel := make([]string, 0, len(u.accounts))
	for k, a := range u.accounts {
		if a.total() != 0 {
			notDel = append(notDel, k)
		}
	}
	return notDel
}

func (u *Unit) start() []string {
	notStart := make([]string, 0, len(u.accounts))
	for k, a := range u.accounts {
		if !a.start() {
			notStart = append(notStart, k)
		}
	}
	return notStart
}

func (u *Unit) stop() []string {
	notStop := make([]string, 0, len(u.accounts))
	for k, a := range u.accounts {
		if !a.stop() {
			notStop = append(notStop, k)
		}
	}
	return notStop
}
