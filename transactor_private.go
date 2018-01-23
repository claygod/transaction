package transactor

// Transactor
// Private
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"sync/atomic"
)

func (t *Transactor) catch() bool {
	if atomic.LoadInt64(&t.hasp) == stateOpen {
		atomic.AddInt64(&t.counter, 1)
		return true
	}
	return false
}
func (t *Transactor) throw() {
	atomic.AddInt64(&t.counter, -1)
}

/*
func (t *Transactor) getAccount(id int64, key string) (*Account, errorCodes) {
	un, ok := t.units.Load(id)
	if !ok {
		t.lgr.New().Context("Msg", errMsgUnitNotExist).Context("Unit", id).Context("Account", id).Context("Method", "getAccount").Write()
		return nil, ErrCodeUnitNotExist
	}
	u := un.(*Unit)
	return u.getAccount(key), Ok
}
*/
func (t *Transactor) getAccount(id int64, key string) (*Account, errorCodes) {
	u := t.storage.getUnit(id)
	if u == nil {
		return nil, ErrCodeUnitNotExist
	}
	return u.getAccount(key), Ok
}
