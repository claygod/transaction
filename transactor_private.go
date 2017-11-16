package transactor

// Transactor
// Private
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"errors"
	//"fmt"
	//"log"
	//"bytes"
	//"os"
	//"runtime"
	//"strconv"
	//"strings"
	//"io/ioutil"
	//"sync"
	"sync/atomic"
)

func (t *Transactor) catch() bool {
	if atomic.LoadInt64(&t.hasp) > 0 {
		atomic.AddInt64(&t.counter, 1)
		return true
	}
	return false
}
func (t *Transactor) throw() {
	atomic.AddInt64(&t.counter, -1)
}

func (t *Transactor) getAccount(id int64, key string) (*Account, errorCodes) {
	// sync.Map begin
	un, ok := t.units.Load(id)
	if !ok {
		t.lgr.New().Context("Msg", errMsgUnitNotExist).Context("Unit", id).Context("Account", id).Context("Method", "getAccount").Write()
		return nil, ErrCodeUnitNotExist
	}
	u := un.(*Unit)
	return u.getAccount(key), Ok
	// sync.Map end
	/*
		t.m.Lock()
		u, ok := t.Units[id]
		t.m.Unlock()
		if !ok {
			t.lgr.New().Context("Msg", errMsgUnitExist).Context("Unit", id).Context("Account", id).Context("Method", "getAccount").Write()
			return nil, ErrCodeUnitExist
		}
		return u.getAccount(key), Ok
	*/
}

//func (t *Transactor) getNEL() []byte {
//	return []byte("")
//}