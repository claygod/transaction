package transactor

// Transactor
// Unit
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"errors"
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

func (u *Unit) Account(key string) *Account {
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

func (u *Unit) List() []string {
	lst := make([]string, 0, len(u.accounts))
	for k, _ := range u.accounts {
		lst = append(lst, k)
	}
	return lst
}

func (u *Unit) Total() map[string]int64 {
	t := make(map[string]int64)
	for k, a := range u.accounts {
		t[k] = a.Total()
	}
	return t
}

func (u *Unit) DelAccount(key string) error {
	a, ok := u.accounts[key]
	if !ok {
		return errors.New("There is no such account")
	}
	if a.Total() != 0 {
		return errors.New("Account is not empty")
	}
	if !a.stop() {
		return errors.New("Account does not stop")
	}
	u.m.Lock()
	delete(u.accounts, key)
	u.m.Unlock()
	return nil
}
func (u *Unit) delAllAccounts() ([]string, error) {
	if notDel := u.del(); len(notDel) != 0 {
		return notDel, errors.New("Accounts are not empty")
	}
	if notStop := u.stop(); len(notStop) != 0 {
		return notStop, errors.New("Do not stop accounts")
	}
	return nil, nil
}

func (u *Unit) del() []string {
	notDel := make([]string, 0, len(u.accounts))
	for k, a := range u.accounts {
		if a.Total() != 0 {
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
