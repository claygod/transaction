package transactor

// Transactor
// API
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
)

type Transactor struct {
	m       sync.Mutex
	counter int64
	hasp    int64
	//units   sync.Map
	lgr     *logger
	storage *Storage
}

// New - create new transactor.
func New() Transactor {
	return Transactor{
		hasp:    stateOpen,
		lgr:     &logger{},
		storage: newStorage(),
	}
}

func (t *Transactor) AddUnit(id int64) errorCodes {
	if !t.catch() {
		go t.lgr.New().Context("Msg", errMsgTransactorNotCatch).Context("Unit", id).Context("Method", "AddUnit").Write()
		return ErrCodeTransactorCatch
	}
	defer t.throw()

	if !t.storage.addUnit(id) { // _, ok := t.units.LoadOrStore(id, newUnit()); ok
		go t.lgr.New().Context("Msg", errMsgUnitExist).Context("Unit", id).Context("Method", "AddUnit").Write()
		return ErrCodeUnitExist
	}
	return Ok
}

func (t *Transactor) DelUnit(id int64) ([]string, errorCodes) {
	if !t.catch() {
		go t.lgr.New().Context("Msg", errMsgTransactorNotCatch).Context("Unit", id).Context("Method", "DelUnit").Write()
		return nil, ErrCodeTransactorCatch
	}
	defer t.throw()

	un, ok := t.storage.delUnit(id) // t.units.Load(id)
	if !ok {
		go t.lgr.New().Context("Msg", errMsgUnitNotExist).Context("Unit", id).Context("Method", "DelUnit").Write()
		return nil, ErrCodeUnitNotExist
	}
	//u := un //.(*Unit)
	if accList, err := un.delAllAccounts(); err != Ok {
		go t.lgr.New().Context("Msg", err).Context("Unit", id).Context("Method", "DelUnit").Write()
		return accList, err
	}
	return nil, Ok
}

func (t *Transactor) TotalUnit(id int64) (map[string]int64, errorCodes) {
	if !t.catch() {
		go t.lgr.New().Context("Msg", errMsgTransactorNotCatch).Context("Unit", id).Context("Method", "TotalUnit").Write()
		return nil, ErrCodeTransactorCatch
	}
	defer t.throw()

	un := t.storage.getUnit(id) // t.units.Load(id)
	if un == nil {
		go t.lgr.New().Context("Msg", errMsgUnitNotExist).Context("Unit", id).Context("Method", "TotalUnit").Write()
		return nil, ErrCodeUnitNotExist
	}
	//u := un //.(*Unit)

	return un.total(), Ok
}

func (t *Transactor) TotalAccount(id int64, key string) (int64, errorCodes) {
	if !t.catch() {
		go t.lgr.New().Context("Msg", errMsgTransactorNotCatch).Context("Unit", id).Context("Account", key).Context("Method", "TotalAccount").Write()
		return permitError, ErrCodeTransactorCatch
	}
	defer t.throw()

	un := t.storage.getUnit(id) // t.units.Load(id)
	if un == nil {
		go t.lgr.New().Context("Msg", errMsgUnitNotExist).Context("Unit", id).Context("Account", key).Context("Method", "TotalAccount").Write()
		return permitError, ErrCodeUnitNotExist
	}
	//u := un.(*Unit)
	return un.getAccount(key).total(), Ok
}

func (t *Transactor) Start() bool {
	for i := trialLimit; i > trialStop; i-- {
		if atomic.LoadInt64(&t.hasp) == stateOpen || atomic.CompareAndSwapInt64(&t.hasp, stateClosed, stateOpen) {
			return true
		}
		runtime.Gosched()
	}
	go t.lgr.New().Context("Msg", errMsgTransactorNotStart).Context("Method", "Start").Write()
	return false
}

func (t *Transactor) Stop() bool {
	for i := trialLimit; i > trialStop; i-- {
		if atomic.LoadInt64(&t.hasp) == stateClosed || (atomic.LoadInt64(&t.counter) == 0 && atomic.CompareAndSwapInt64(&t.hasp, stateOpen, stateClosed)) {
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
		un := t.storage.getUnit(id)
		if un == nil {
			t.storage.addUnit(id)
			un = t.storage.getUnit(id)
		}

		un.accounts[string(a[2])] = newAccount(balance)
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
	for i := uint64(0); i < storageNumber; i++ {
		for id, u := range t.storage.data[i].data {
			for key, balance := range u.totalUnsave() {
				buf.Write([]byte(fmt.Sprintf("%d%s%d%s%s%s", id, separatorSymbol, balance, separatorSymbol, key, endLineSymbol)))
			}
		}
	}

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

func (t *Transactor) Unsafe(reqs []*Request) errorCodes {
	tn := &Transaction{
		tr:   t,
		up:   make([]*Request, 0, 0),
		reqs: reqs,
	}
	return tn.exeTransaction()
}
