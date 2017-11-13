package transactor

// Transactor
// API
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"errors"
	"fmt"
	//"log"
	"bytes"
	"os"
	"runtime"
	"strconv"
	//"strings"
	"io/ioutil"
	"sync"
	"sync/atomic"
)

type Transactor struct {
	m       sync.Mutex
	counter int64
	hasp    int64
	Units   map[int64]*Unit
	lgr     *logger
}

// New - create new transactor.
func New() Transactor {
	t := Transactor{hasp: 0, Units: make(map[int64]*Unit), lgr: &logger{}}

	//t.lgr.New().Context("TEST", "LOG").Context("Type", ErrLevelError).
	//	Context("Msg", ErrMsgUnitExist).Context("Unit", 1234242343).Write()
	return t
}

func (t *Transactor) AddUnit(id int64) errorCodes {
	if !t.catch() {
		return ErrCodeTransactorCatch
	}
	defer t.throw()
	_, ok := t.Units[id]
	if !ok {
		t.m.Lock()
		_, ok = t.Units[id]
		if !ok {
			t.Units[id] = newUnit()
			t.m.Unlock()
			return ErrOk
		}
		t.m.Unlock()
	}
	t.lgr.New().Context("Msg", ErrMsgUnitExist).Context("Unit", id).Context("Method", "AddUnit").Write()
	return ErrCodeUnitExist
}

/*
func (t *Transactor) GetUnit(id int64) (*Unit, errorCodes) {
	u, ok := t.Units[id]
	if !ok {
		t.lgr.New().Context("Msg", ErrMsgUnitExist).Context("Unit", id).Context("Method", "GetUnit").Write()
		return nil, ErrCodeUnitExist
	}
	return u, ErrOk
}
*/

func (t *Transactor) Total() (map[int64]map[string]int64, errorCodes) {
	if !t.catch() {
		return nil, ErrCodeTransactorCatch
	}
	defer t.throw()
	ttl := make(map[int64]map[string]int64)
	for k, u := range t.Units {
		ttl[k] = u.total()
	}
	return ttl, ErrOk
}

func (t *Transactor) TotalUnit(id int64) (map[string]int64, errorCodes) {
	u, ok := t.Units[id]
	if !ok {
		t.lgr.New().Context("Msg", ErrMsgUnitNotExist).Context("Unit", id).Context("Method", "TotalUnit").Write()
		return nil, ErrCodeUnitNotExist
	}
	return u.total(), ErrOk
}

func (t *Transactor) TotalAccount(id int64, key string) (int64, errorCodes) {
	u, ok := t.Units[id]
	if !ok {
		t.lgr.New().Context("Msg", ErrMsgUnitNotExist).Context("Unit", id).Context("Method", "TotalAccount").Write()
		return -1, ErrCodeUnitNotExist
	}
	return u.getAccount(key).total(), ErrOk
}

func (t *Transactor) Start() bool {
	for i := trialLimit; i > trialStop; i-- {
		if atomic.LoadInt64(&t.hasp) == 1 || atomic.CompareAndSwapInt64(&t.hasp, 0, 1) {
			return true
		}
		runtime.Gosched()
	}
	return false
}

func (t *Transactor) Stop() bool {
	for i := trialLimit; i > trialStop; i-- {
		if (atomic.LoadInt64(&t.hasp) == 0 || atomic.CompareAndSwapInt64(&t.hasp, 1, 0)) && atomic.LoadInt64(&t.counter) == 0 {
			return true
		}
		runtime.Gosched()
	}
	return false
}

func (t *Transactor) Load(path string) errorCodes {
	hasp := atomic.LoadInt64(&t.hasp)
	if hasp == 1 && !t.Stop() {
		return ErrCodeTransactorStop
	}

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return ErrCodeLoadReadFile
	}
	for _, str := range bytes.Split(bs, []byte("\r\n")) {
		a := bytes.Split(str, []byte(";"))
		if len(a) != 3 {
			continue
		}
		id, err := strconv.ParseInt(string(a[0]), 10, 64)
		if err != nil {
			return ErrCodeLoadStrToInt64
		}
		balance, err := strconv.ParseInt(string(a[1]), 10, 64)
		if err != nil {
			return ErrCodeLoadStrToInt64
		}
		u, ok := t.Units[id]
		if !ok {
			u = newUnit()
			t.Units[id] = u
		}
		u.accounts[string(a[2])] = newAccount(balance)
	}
	if hasp == 1 && !t.Start() {
		return ErrCodeTransactorStart
	}
	return ErrOk
}

func (t *Transactor) Save(path string) errorCodes {
	hasp := atomic.LoadInt64(&t.hasp)
	if hasp == 1 && !t.Stop() {
		return ErrCodeTransactorStop
	}

	var buf bytes.Buffer
	for id, u := range t.Units {
		for key, a := range u.accounts {
			buf.Write([]byte(fmt.Sprintf("%d;%d;%s\r\n", id, a.balance, key)))
		}
	}
	//file, err := os.Create(path)
	//if err != nil {
	//	file.Close()
	//	return ErrCodeSaveCreateFile //fmt.Print("=!=", err, "=!=")
	//}
	//file.Write(buf.Bytes())
	//file.Close()
	if ioutil.WriteFile(path, buf.Bytes(), os.FileMode(0777)) != nil {
		return ErrCodeSaveCreateFile
	}
	if hasp == 1 && !t.Start() {
		return ErrCodeTransactorStart
	}
	return ErrOk
}

func (t *Transactor) DelUnit(id int64) ([]string, errorCodes) {
	if !t.catch() {
		return nil, ErrCodeTransactorCatch
	}
	defer t.throw()
	if u, ok := t.Units[id]; ok {
		if accList, err := u.delAllAccounts(); err != ErrOk {
			t.lgr.New().Context("Msg", err).Context("Unit", id).Context("Method", "DelUnit").Write()
			return accList, err
		}
	}

	return nil, ErrOk
}

func (t *Transactor) DelAccount(id int64, key string) errorCodes {
	if !t.catch() {
		return ErrCodeTransactorCatch
	}
	defer t.throw()
	if u, ok := t.Units[id]; ok {
		return u.delAccount(key)
	}
	t.lgr.New().Context("Msg", ErrMsgUnitNotExist).Context("Unit", id).Context("Method", "DelAccount").Write()
	return ErrCodeUnitNotExist
}

func (t *Transactor) Begin() *Transaction {
	return newTransaction(t)
}
