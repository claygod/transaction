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

type Customer struct {
	m        sync.Mutex
	accounts map[string]*Account
}

// newCustomer - create new account.
func newCustomer() *Customer {
	k := &Customer{accounts: make(map[string]*Account)}
	return k
}

// Account - get an account link.
// If there is no such account, it will be created.
func (c *Customer) Account(num string) *Account {
	a, ok := c.accounts[num]
	if !ok {
		c.m.Lock()
		a, ok = c.accounts[num]
		if !ok {
			a = newAccount(0)
			c.accounts[num] = a
		}
		c.m.Unlock()
	}
	return a
}

func (c *Customer) delAllAccounts(num string) map[string]*Account {
	aList := make(map[string]*Account)
	c.m.Lock()
	for key, a := range c.accounts {
		if c.delAccountNoLockUnsafe(key) != nil {
			aList[key] = a
		} else {
			delete(c.accounts, key)
		}
	}
	c.m.Unlock()
	return aList
}

func (c *Customer) DelAccount(num string) (int64, int64, error) {
	_, ok := c.accounts[num]
	if ok {
		c.m.Lock()
		defer c.m.Unlock()
		a, ok := c.accounts[num]
		if ok {
			if a.balance == 0 && a.debt == 0 {
				delete(c.accounts, num)
				return 0, 0, nil
			}
			return a.balance, a.debt, errors.New("Account is not zero.")
		}
	}
	return -1, -1, errors.New("There is no such account")
}

func (c *Customer) delAccountNoLockUnsafe(num string) error {
	a := c.accounts[num]
	if a.debt == 0 && a.balance == 0 {
		delete(c.accounts, num)
		return nil
	}
	return errors.New("Account is not zero")
}

/*
func (c *customer) createAccount(num string, amount int64) error {
	if _, ok := c.accounts[num]; ok {
		return errors.New("This account already exists")
	}
	c.accounts[num] = newAccount(amount)
	return nil
}
*/
