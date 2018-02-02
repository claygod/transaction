package transaction

// Core
// Private
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"sync/atomic"
)

/*
catch - obtaining a permit for an operation.
If we have an open status, the counter is incremented
and the `true` returns, otherwise the `false` is returned.
*/
func (c *Core) catch() bool {
	if atomic.LoadInt64(&c.hasp) == stateOpen {
		atomic.AddInt64(&c.counter, 1)
		return true
	}
	return false
}

/*
throw - permission to conduct an operation is no longer required.
Decrement of the counter.
*/
func (c *Core) throw() {
	atomic.AddInt64(&c.counter, -1)
}

/*
getAccount - get account by ID and string key.
If a unit with such an ID exists, then the account and code `Ok` is returned.

Returned codes:
	ErrCodeUnitNotExist // a unit with such an ID does not exist
	Ok
*/
func (c *Core) getAccount(id int64, key string) (*account, errorCodes) {
	u := c.storage.getUnit(id)
	if u == nil {
		return nil, ErrCodeUnitNotExist
	}
	return u.getAccount(key), Ok
}
