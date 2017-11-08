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
const permitError int64 = -9223372036854775806

type Transactor struct {
	m     sync.Mutex
	Units map[int64]*Unit
}

// New - create new transactor.
func New() Transactor {
	k := Transactor{Units: make(map[int64]*Unit)}
	return k
}

func (t *Transactor) AddUnit(id int64) error {
	_, ok := t.Units[id]
	if !ok {
		t.m.Lock()
		defer t.m.Unlock()
		_, ok = t.Units[id]
		if !ok {
			t.Units[id] = newUnit()
			return nil
		}
	}
	return errors.New("This unit already exists")
}

func (t *Transactor) GetUnit(id int64) *Unit {
	u, ok := t.Units[id]
	if !ok {
		return nil //errors.New("This unit does not exist")
	}
	return u
}

func (t *Transactor) DelUnit(id int64) ([]string, error) {
	if u, ok := t.Units[id]; ok {
		if accList, err := u.delAllAccounts(); err != nil {
			return accList, err
		}
	}
	return nil, nil
}

func (t *Transactor) getAccount(id int64, key string) (*Account, error) {
	u, ok := t.Units[id]
	if !ok {
		return nil, errors.New("This unit already exists")
	}
	return u.Account(key), nil
}

func (t *Transactor) Begin() *Transaction {
	return newTransaction(t)
}

func (t *Transactor) Total() map[int64]map[string]int64 {
	ttl := make(map[int64]map[string]int64)
	for k, u := range t.Units {
		ttl[k] = u.Total()
	}
	return ttl
}
