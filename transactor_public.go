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
	units   sync.Map
	lgr     *logger
}

// New - create new transactor.
func New() Transactor {
	return Transactor{hasp: stateOpen, lgr: &logger{}}
}

func (t *Transactor) AddUnit(id int64) errorCodes {
	if !t.catch() {
		go t.lgr.New().Context("Msg", errMsgTransactorNotCatch).Context("Unit", id).Context("Method", "AddUnit").Write()
		return ErrCodeTransactorCatch
	}
	defer t.throw()
	// sync.Map begin
	if _, ok := t.units.LoadOrStore(id, newUnit()); ok {
		go t.lgr.New().Context("Msg", errMsgUnitExist).Context("Unit", id).Context("Method", "AddUnit").Write()
		return ErrCodeUnitExist
	}
	return Ok
	// sync.Map end
	/*
		_, ok := t.Units[id]
		if !ok {
			t.m.Lock()
			_, ok = t.Units[id]
			if !ok {
				t.Units[id] = newUnit()
				t.m.Unlock()
				return Ok
			}
			t.m.Unlock()
		}
		go t.lgr.New().Context("Msg", errMsgUnitExist).Context("Unit", id).Context("Method", "AddUnit").Write()
		return ErrCodeUnitExist
	*/
}

func (t *Transactor) DelUnit(id int64) ([]string, errorCodes) {
	if !t.catch() {
		go t.lgr.New().Context("Msg", errMsgTransactorNotCatch).Context("Unit", id).Context("Method", "DelUnit").Write()
		return nil, ErrCodeTransactorCatch
	}
	defer t.throw()
	// sync.Map begin
	un, ok := t.units.Load(id)
	if !ok {
		go t.lgr.New().Context("Msg", errMsgUnitNotExist).Context("Unit", id).Context("Method", "DelUnit").Write()
		return nil, ErrCodeUnitNotExist
	}
	u := un.(*Unit)
	if accList, err := u.delAllAccounts(); err != Ok {
		go t.lgr.New().Context("Msg", err).Context("Unit", id).Context("Method", "DelUnit").Write()
		return accList, err
	}
	return nil, Ok
	// sync.Map end
	/*
		t.m.Lock()
		if u, ok := t.Units[id]; ok {
			if accList, err := u.delAllAccounts(); err != Ok {
				t.m.Unlock()
				go t.lgr.New().Context("Msg", err).Context("Unit", id).Context("Method", "DelUnit").Write()
				return accList, err
			}
		}
		t.m.Unlock()
		return nil, Ok
	*/
}

func (t *Transactor) Total() (map[int64]map[string]int64, errorCodes) {
	if !t.catch() {
		go t.lgr.New().Context("Msg", errMsgTransactorNotCatch).Context("Method", "Total").Write()
		return nil, ErrCodeTransactorCatch
	}
	defer t.throw()

	ttl := make(map[int64]map[string]int64)

	t.units.Range(func(id, u interface{}) bool {
		id2 := id.(int64)
		u2 := u.(*Unit)
		ttl[id2] = u2.total()
		return true // if false, Range stops
	})

	//for k, u := range t.Units {
	//	ttl[k] = u.total()
	//}
	return ttl, Ok
}

func (t *Transactor) TotalUnit(id int64) (map[string]int64, errorCodes) {
	if !t.catch() {
		go t.lgr.New().Context("Msg", errMsgTransactorNotCatch).Context("Unit", id).Context("Method", "TotalUnit").Write()
		return nil, ErrCodeTransactorCatch
	}
	defer t.throw()

	// sync.Map begin
	un, ok := t.units.Load(id)
	if !ok {
		go t.lgr.New().Context("Msg", errMsgUnitNotExist).Context("Unit", id).Context("Method", "TotalUnit").Write()
		return nil, ErrCodeUnitNotExist
	}
	u := un.(*Unit)
	// sync.Map end
	/*
		u, ok := t.Units[id]
		if !ok {
			go t.lgr.New().Context("Msg", errMsgUnitNotExist).Context("Unit", id).Context("Method", "TotalUnit").Write()
			return nil, ErrCodeUnitNotExist
		}
	*/
	return u.total(), Ok
}

func (t *Transactor) TotalAccount(id int64, key string) (int64, errorCodes) {
	if !t.catch() {
		go t.lgr.New().Context("Msg", errMsgTransactorNotCatch).Context("Unit", id).Context("Account", key).Context("Method", "TotalAccount").Write()
		return permitError, ErrCodeTransactorCatch
	}
	defer t.throw()

	// sync.Map begin
	un, ok := t.units.Load(id)
	if !ok {
		go t.lgr.New().Context("Msg", errMsgUnitNotExist).Context("Unit", id).Context("Account", key).Context("Method", "TotalAccount").Write()
		return permitError, ErrCodeUnitNotExist
	}
	u := un.(*Unit)
	return u.getAccount(key).total(), Ok
	// sync.Map end
	/*
		u, ok := t.Units[id]
		if !ok {
			go t.lgr.New().Context("Msg", errMsgUnitNotExist).Context("Unit", id).Context("Method", "TotalAccount").Write()
			return permitError, ErrCodeUnitNotExist
		}
		return u.getAccount(key).total(), Ok
	*/
}

func (t *Transactor) Start() bool {
	for i := trialLimit; i > trialStop; i-- {
		if atomic.LoadInt64(&t.hasp) == stateClosed || atomic.CompareAndSwapInt64(&t.hasp, stateOpen, stateClosed) {
			return true
		}
		runtime.Gosched()
	}
	go t.lgr.New().Context("Msg", errMsgTransactorNotStart).Context("Method", "Start").Write()
	return false
}

func (t *Transactor) Stop() bool {
	for i := trialLimit; i > trialStop; i-- {
		if (atomic.LoadInt64(&t.hasp) == stateOpen || atomic.CompareAndSwapInt64(&t.hasp, stateClosed, stateOpen)) && atomic.LoadInt64(&t.counter) == 0 {
			return true
		}
		runtime.Gosched()
	}
	go t.lgr.New().Context("Msg", errMsgTransactorNotStop).Context("Method", "Stop").Write()
	return false
}

func (t *Transactor) Load(path string) errorCodes {
	hasp := atomic.LoadInt64(&t.hasp)
	if hasp == stateClosed && !t.Stop() {
		go t.lgr.New().Context("Msg", errMsgTransactorNotStop).Context("Method", "Load").Write()
		return ErrCodeTransactorStop
	}

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		go t.lgr.New().Context("Msg", errMsgTransactorNotReadFile).Context("Path", path).Context("Method", "Load").Write()
		return ErrCodeLoadReadFile
	}
	endLine := []byte(endLineSymbol)
	separator := []byte(separatorSymbol)
	for _, str := range bytes.Split(bs, endLine) {
		a := bytes.Split(str, separator)
		if len(a) != 3 {
			continue
		}
		id, err := strconv.ParseInt(string(a[0]), 10, 64)
		if err != nil {
			go t.lgr.New().Context("Msg", errMsgTransactorParseString).Context("Path", path).Context("String", str).Context("Method", "Load").Write()
			return ErrCodeLoadStrToInt64
		}
		balance, err := strconv.ParseInt(string(a[1]), 10, 64)
		if err != nil {
			go t.lgr.New().Context("Msg", errMsgTransactorParseString).Context("Path", path).Context("String", str).Context("Method", "Load").Write()
			return ErrCodeLoadStrToInt64
		}
		// sync.Map begin
		un, _ := t.units.LoadOrStore(id, newUnit())
		u := un.(*Unit)
		// sync.Map end

		//u, ok := t.Units[id]
		//if !ok {
		//	u = newUnit()
		//	t.Units[id] = u
		//}
		u.accounts[string(a[2])] = newAccount(balance)
	}
	if hasp == stateClosed && !t.Start() {
		go t.lgr.New().Context("Msg", errMsgTransactorNotStart).Context("Method", "Load").Write()
		return ErrCodeTransactorStart
	}
	return Ok
}

func (t *Transactor) Save(path string) errorCodes {
	hasp := atomic.LoadInt64(&t.hasp)
	if hasp == stateClosed && !t.Stop() {
		go t.lgr.New().Context("Msg", errMsgTransactorNotStop).Context("Method", "Save").Write()
		return ErrCodeTransactorStop
	}

	var buf bytes.Buffer
	// sync.Map begin
	t.units.Range(func(id, u interface{}) bool {
		id2 := id.(int64)
		u2 := u.(*Unit)
		for key, a := range u2.accounts {
			buf.Write([]byte(fmt.Sprintf("%d%s%d%s%s%s", id2, separatorSymbol, a.balance, separatorSymbol, key, endLineSymbol)))
		}
		return true // if false, Range stops
	})
	// sync.Map end
	/*
		for id, u := range t.Units {
			for key, a := range u.accounts {
				buf.Write([]byte(fmt.Sprintf("%d%s%d%s%s%s", id, separatorSymbol, a.balance, separatorSymbol, key, endLineSymbol)))
			}
		}
	*/
	if ioutil.WriteFile(path, buf.Bytes(), os.FileMode(0777)) != nil {
		go t.lgr.New().Context("Msg", errMsgTransactorNotCreateFile).Context("Path", path).Context("Method", "Save").Write()
		return ErrCodeSaveCreateFile
	}
	if hasp == stateClosed && !t.Start() {
		go t.lgr.New().Context("Msg", errMsgTransactorNotStart).Context("Method", "Save").Write()
		return ErrCodeTransactorStart
	}
	return Ok
}

func (t *Transactor) Begin() *Transaction {
	return newTransaction(t)
}
