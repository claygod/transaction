package transactor

// Core
// Private
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"sync/atomic"
)

func (c *Core) catch() bool {
	if atomic.LoadInt64(&c.hasp) == stateOpen {
		atomic.AddInt64(&c.counter, 1)
		return true
	}
	return false
}
func (c *Core) throw() {
	atomic.AddInt64(&c.counter, -1)
}

/*
func (t *Core) getAccount(id int64, key string) (*Account, errorCodes) {
	un, ok := t.units.Load(id)
	if !ok {
		t.lgr.New().Context("Msg", errMsgUnitNotExist).Context("Unit", id).Context("Account", id).Context("Method", "getAccount").Write()
		return nil, ErrCodeUnitNotExist
	}
	u := un.(*Unit)
	return u.getAccount(key), Ok
}
*/
func (c *Core) getAccount(id int64, key string) (*Account, errorCodes) {
	u := c.storage.getUnit(id)
	if u == nil {
		return nil, ErrCodeUnitNotExist
	}
	return u.getAccount(key), Ok
}
