package transactor

// Transactor
// API
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"errors"
	//"fmt"
	//"log"
	//"runtime"
	"sync"
	//"sync/atomic"
)

const countNodes int = 65536
const trialLimit int = 20000000
const trialStop int = 64

type Transactor struct {
	m         sync.Mutex
	customers map[int64]*Customer
}

// New - create new transactor.
func New() Transactor {
	k := Transactor{customers: make(map[int64]*Customer)}
	return k
}

func (t *Transactor) AddCustomer(id int64) error {
	_, ok := t.customers[id]
	if !ok {
		t.m.Lock()
		defer t.m.Unlock()
		_, ok = t.customers[id]
		if !ok {
			t.customers[id] = newCustomer()
			return nil
		}
	}
	return errors.New("This customer already exists")
}

func (t *Transactor) getAccount(id int64, key string) *Account {
	c, ok := t.customers[id]
	if !ok {
		return nil
	}
	return c.Account(key)
}

func (t *Transactor) Begin() *Transaction {
	return newTransaction(t)
}

func (t *Transactor) Customer(cid int64) *Customer {
	c, ok := t.customers[cid]
	if !ok {
		return nil //errors.New("This customer does not exist")
	}
	return c
}

func (t *Transactor) AccountStore(cid int64, key string) (int64, int64, error) {
	c, ok := t.customers[cid]
	if !ok {
		return -2, -2, errors.New("There is no such customer")
	}
	return c.AccountStore(key)
}

func (t *Transactor) DelCustomer(cid int64) error {
	_, ok := t.customers[cid]
	if !ok {
		return errors.New("This customer does not exist")
	}
	//
	return nil
}
