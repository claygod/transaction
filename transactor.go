package transactor

// Transactor
// API
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"errors"
	//"fmt"
	"log"
	//"runtime"
	"sync"
	//"sync/atomic"
)

type Transactor struct {
	m      sync.Mutex
	Units  map[int64]*Unit
	lgr    *logger
	writer log.Logger
}

// New - create new transactor.
func New() Transactor {
	t := Transactor{Units: make(map[int64]*Unit), lgr: &logger{}}

	t.lgr.New().Context("TEST", "LOG").Context("Type", ErrLevelError).
		Context("Msg", ErrMsgUnitExist).Context("Unit", 1234242343).Write()
	return t
}

func (t *Transactor) AddUnit(id int64) errorCodes {
	_, ok := t.Units[id]
	if !ok {
		t.m.Lock()
		defer t.m.Unlock()
		_, ok = t.Units[id]
		if !ok {
			t.Units[id] = newUnit()
			return ErrOk
		}
	}
	t.lgr.New().Context("Msg", ErrMsgUnitExist).Context("Unit", id).Context("Method", "AddUnit").Write()
	return ErrCodeUnitExist
}

func (t *Transactor) GetUnit(id int64) (*Unit, errorCodes) {
	u, ok := t.Units[id]
	if !ok {
		t.lgr.New().Context("Msg", ErrMsgUnitExist).Context("Unit", id).Context("Method", "GetUnit").Write()
		return nil, ErrCodeUnitExist
	}
	return u, ErrOk
}

func (t *Transactor) DelUnit(id int64) ([]string, errorCodes) {
	if u, ok := t.Units[id]; ok {
		if accList, err := u.delAllAccounts(); err != ErrOk {
			t.lgr.New().Context("Msg", err).Context("Unit", id).Context("Method", "DelUnit").Write()
			return accList, err
		}
	}
	return nil, ErrOk
}

func (t *Transactor) getAccount(id int64, key string) (*Account, errorCodes) {
	u, ok := t.Units[id]
	if !ok {
		t.lgr.New().Context("Msg", ErrMsgUnitExist).Context("Unit", id).Context("Account", id).Context("Method", "getAccount").Write()
		return nil, ErrCodeUnitExist
	}
	return u.Account(key), ErrOk
}

func (t *Transactor) Begin() *Transaction {
	return newTransaction(t)
}

func (t *Transactor) Total() map[int64]map[string]int64 {
	log.Print(1)
	ttl := make(map[int64]map[string]int64)
	for k, u := range t.Units {
		ttl[k] = u.Total()
	}
	return ttl
}
