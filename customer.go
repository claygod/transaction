package transaction

// Transaction
// Customer
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"errors"
	//"fmt"
	//"runtime"
	"sync"
	//"sync/atomic"
)

type customer struct {
	m        sync.Mutex
	accounts map[string]*account
}

// newCustomer - create new account.
func newCustomer(amount int64) customer {
	k := customer{accounts: make(map[string]*account)}
	return k
}

func (c *customer) account(acnt string) *account {
	a, ok := c.accounts[acnt]
	if !ok {
		a = newAccount(0)
		c.accounts[acnt] = a
	}
	return a
}

func (c *customer) delAccount(acnt string) error {
	_, ok := c.accounts[acnt]
	if ok {
		c.m.Lock()
		delete(c.accounts, acnt)
		c.m.Unlock()
		return nil
	}
	return errors.New("There is no such account")
}

/*
func (c *customer) createAccount(acnt string, amount int64) error {
	if _, ok := c.accounts[acnt]; ok {
		return errors.New("This account already exists")
	}
	c.accounts[acnt] = newAccount(amount)
	return nil
}
*/
