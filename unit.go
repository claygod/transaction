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
	m        sync.Mutex
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
		u.m.Lock()
		a, ok = u.accounts[key]
		if !ok {
			a = newAccount(0)
			u.accounts[key] = a
		}
		u.m.Unlock()
	}
	return a
}

/*
func (u *Unit) List() []string {
	lst := make([]string, 0, len(u.accounts))
	for k, _ := range u.accounts {
		lst = append(lst, k)
	}
	return lst
}
*/
func (u *Unit) total() map[string]int64 {
	t := make(map[string]int64)
	for k, a := range u.accounts {
		t[k] = a.total()
	}
	return t
}

func (u *Unit) delAccount(key string) errorCodes {
	a, ok := u.accounts[key]
	if !ok {
		return ErrCodeAccountExist
	}
	if a.total() != 0 {
		return ErrCodeAccountNotEmpty
	}
	if !a.stop() {
		return ErrCodeAccountNotStop
	}
	u.m.Lock()
	delete(u.accounts, key)
	u.m.Unlock()
	return Ok
}
func (u *Unit) delAllAccounts() ([]string, errorCodes) {
	if notDel := u.del(); len(notDel) != 0 {
		return notDel, ErrCodeAccountNotEmpty
	}
	if notStop := u.stop(); len(notStop) != 0 {
		return notStop, ErrCodeAccountNotStop
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

func (u *Unit) stop() []string {
	notStop := make([]string, 0, len(u.accounts))
	for k, a := range u.accounts {
		if !a.stop() {
			notStop = append(notStop, k)
		}
	}
	return notStop
}
