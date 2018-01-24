package transactor

// Core
// Public
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

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

// Core - root application structure
type Core struct {
	m       sync.Mutex
	counter int64
	hasp    int64
	lgr     *logger
	storage *Storage
}

// New - create new core
func New() Core {
	return Core{
		hasp:    stateOpen,
		lgr:     &logger{},
		storage: newStorage(),
	}
}

/*
AddUnit - adding a new unit.
Two units with the same identifier can not exist.

Returned codes:

	ErrCodeCoreCatch - not obtained permission
	ErrCodeUnitExist - such a unit already exists
	Ok
*/
func (c *Core) AddUnit(id int64) errorCodes {
	if !c.catch() {
		go c.lgr.New().Context("Msg", errMsgCoreNotCatch).Context("Unit", id).Context("Method", "AddUnit").Write()
		return ErrCodeCoreCatch
	}
	defer c.throw()

	if !c.storage.addUnit(id) {
		go c.lgr.New().Context("Msg", errMsgUnitExist).Context("Unit", id).Context("Method", "AddUnit").Write()
		return ErrCodeUnitExist
	}
	return Ok
}

// DelUnit - deletion of a unit
// The unit will be deleted only when all its accounts are stopped and deleted.
// In case of an error, a list of not deleted accounts is returned.
// Returned codes:
// - ErrCodeCoreCatch - not obtained permission
// - ErrCodeUnitExist - there is no such unit
// - ErrCodeAccountNotStop - accounts failed to stop
// - ErrCodeUnitNotEmpty - accounts are not zero
// - Ok
func (c *Core) DelUnit(id int64) ([]string, errorCodes) {
	if !c.catch() {
		go c.lgr.New().Context("Msg", errMsgCoreNotCatch).Context("Unit", id).Context("Method", "DelUnit").Write()
		return nil, ErrCodeCoreCatch
	}
	defer c.throw()

	un, ok := c.storage.delUnit(id)
	if !ok {
		go c.lgr.New().Context("Msg", errMsgUnitNotExist).Context("Unit", id).Context("Method", "DelUnit").Write()
		return nil, ErrCodeUnitNotExist
	}

	if accList, err := un.delAllAccounts(); err != Ok {
		go c.lgr.New().Context("Msg", err).Context("Unit", id).Context("Method", "DelUnit").Write()
		return accList, err
	}
	return nil, Ok
}

// TotalUnit - statement on all accounts of the unit
// The ID-balance array is returned.
// Returned codes:
// - ErrCodeCoreCatch - not obtained permission
// - ErrCodeUnitExist - there is no such unit
// - Ok
func (c *Core) TotalUnit(id int64) (map[string]int64, errorCodes) {
	if !c.catch() {
		go c.lgr.New().Context("Msg", errMsgCoreNotCatch).Context("Unit", id).Context("Method", "TotalUnit").Write()
		return nil, ErrCodeCoreCatch
	}
	defer c.throw()

	un := c.storage.getUnit(id)
	if un == nil {
		go c.lgr.New().Context("Msg", errMsgUnitNotExist).Context("Unit", id).Context("Method", "TotalUnit").Write()
		return nil, ErrCodeUnitNotExist
	}

	return un.total(), Ok
}

// TotalAccount - account balance
// If an account has not been created before,
// it will be created with a zero balance.
// Returned codes:
// - ErrCodeCoreCatch - not obtained permission
// - ErrCodeUnitExist - there is no such unit
// - Ok
func (c *Core) TotalAccount(id int64, key string) (int64, errorCodes) {
	if !c.catch() {
		go c.lgr.New().Context("Msg", errMsgCoreNotCatch).Context("Unit", id).Context("Account", key).Context("Method", "TotalAccount").Write()
		return permitError, ErrCodeCoreCatch
	}
	defer c.throw()

	un := c.storage.getUnit(id)
	if un == nil {
		go c.lgr.New().Context("Msg", errMsgUnitNotExist).Context("Unit", id).Context("Account", key).Context("Method", "TotalAccount").Write()
		return permitError, ErrCodeUnitNotExist
	}
	return un.getAccount(key).total(), Ok
}

// Start -
func (c *Core) Start() bool {
	for i := trialLimit; i > trialStop; i-- {
		if atomic.LoadInt64(&c.hasp) == stateOpen || atomic.CompareAndSwapInt64(&c.hasp, stateClosed, stateOpen) {
			return true
		}
		runtime.Gosched()
	}
	go c.lgr.New().Context("Msg", errMsgCoreNotStart).Context("Method", "Start").Write()
	return false
}

func (c *Core) Stop() bool {
	for i := trialLimit; i > trialStop; i-- {
		if atomic.LoadInt64(&c.hasp) == stateClosed || (atomic.LoadInt64(&c.counter) == 0 && atomic.CompareAndSwapInt64(&c.hasp, stateOpen, stateClosed)) {
			return true
		}
		runtime.Gosched()
	}
	go c.lgr.New().Context("Msg", errMsgCoreNotStop).Context("Method", "Stop").Write()
	return false
}

func (c *Core) Load(path string) errorCodes {
	hasp := atomic.LoadInt64(&c.hasp)
	if hasp == stateClosed && !c.Stop() {
		go c.lgr.New().Context("Msg", errMsgCoreNotStop).Context("Method", "Load").Write()
		return ErrCodeCoreStop
	}

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		go c.lgr.New().Context("Msg", errMsgCoreNotReadFile).Context("Path", path).Context("Method", "Load").Write()
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
			go c.lgr.New().Context("Msg", errMsgCoreParseString).Context("Path", path).Context("String", str).Context("Method", "Load").Write()
			return ErrCodeLoadStrToInt64
		}
		balance, err := strconv.ParseInt(string(a[1]), 10, 64)
		if err != nil {
			go c.lgr.New().Context("Msg", errMsgCoreParseString).Context("Path", path).Context("String", str).Context("Method", "Load").Write()
			return ErrCodeLoadStrToInt64
		}
		un := c.storage.getUnit(id)
		if un == nil {
			c.storage.addUnit(id)
			un = c.storage.getUnit(id)
		}

		un.accounts[string(a[2])] = newAccount(balance)
	}
	if hasp == stateClosed && !c.Start() {
		go c.lgr.New().Context("Msg", errMsgCoreNotStart).Context("Method", "Load").Write()
		return ErrCodeCoreStart
	}
	return Ok
}

func (c *Core) Save(path string) errorCodes {
	hasp := atomic.LoadInt64(&c.hasp)
	if hasp == stateClosed && !c.Stop() {
		go c.lgr.New().Context("Msg", errMsgCoreNotStop).Context("Method", "Save").Write()
		return ErrCodeCoreStop
	}

	var buf bytes.Buffer
	for i := uint64(0); i < storageNumber; i++ {
		for id, u := range c.storage.data[i].data {
			for key, balance := range u.totalUnsave() {
				buf.Write([]byte(fmt.Sprintf("%d%s%d%s%s%s", id, separatorSymbol, balance, separatorSymbol, key, endLineSymbol)))
			}
		}
	}

	if ioutil.WriteFile(path, buf.Bytes(), os.FileMode(0777)) != nil {
		go c.lgr.New().Context("Msg", errMsgCoreNotCreateFile).Context("Path", path).Context("Method", "Save").Write()
		return ErrCodeSaveCreateFile
	}
	if hasp == stateClosed && !c.Start() {
		go c.lgr.New().Context("Msg", errMsgCoreNotStart).Context("Method", "Save").Write()
		return ErrCodeCoreStart
	}
	return Ok
}

func (c *Core) Begin() *Transaction {
	return newTransaction(c)
}

/*
func (t *Core) Unsafe(reqs []*Request) errorCodes {
	tn := &Transaction{
		tr:   t,
		up:   make([]*Request, 0, 0),
		reqs: reqs,
	}
	return tn.exeTransaction()
}
*/
