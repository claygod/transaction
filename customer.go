package transaction

// Transaction
// Customer
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"errors"
	//"log"
	"sync"
)

type Customer struct {
	m        sync.Mutex
	accounts sync.Map
}

// newCustomer - create new account.
func newCustomer() *Customer {
	k := &Customer{}
	return k
}

// Account - get an account link.
// If there is no such account, it will be created.
func (c *Customer) Account(key string) *Account {
	a, _ := c.accounts.LoadOrStore(key, newAccount(0))
	return a.(*Account)
}

func (c *Customer) catchAccount(key string) *Account {
	a, _ := c.accounts.LoadOrStore(key, newAccount(0))
	acc := a.(*Account)
	if acc.catch() {
		return a.(*Account)
	}
	return nil
}

func (c *Customer) AccountStore(key string) (int64, int64, error) {
	a, ok := c.accounts.Load(key)
	if !ok {
		return -1, -1, errors.New("There is no such account")
	}
	acc := a.(*Account)
	return acc.balance, acc.debt, nil
}

func (c *Customer) DelAccount(key string) (int64, int64, error) {
	a, ok := c.accounts.Load(key)

	if !ok {
		return -1, -1, errors.New("There is no such account")
	}
	acc := a.(*Account)
	b, d := acc.state()
	if b == 0 && d == 0 {
		c.accounts.Delete(key)
		return 0, 0, nil
	}
	return b, d, errors.New("Account is not zero.")
}

func (c *Customer) Store() map[string][2]int64 {
	nc := make(map[string][2]int64)
	var b, d int64
	acc := newAccount(0)
	c.accounts.Range(func(k, v interface{}) bool {
		acc = v.(*Account)
		b, d = acc.state()
		nc[k.(string)] = [2]int64{b, d}
		return true
	})
	return nc
}
