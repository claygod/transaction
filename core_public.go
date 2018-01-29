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

/*
Core - root application structure
*/
type Core struct {
	m       sync.Mutex
	counter int64
	hasp    int64
	lgr     *logger
	storage *Storage
}

/*
New - create new core
*/
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
	ErrCodeCoreCatch // not obtained permission
	ErrCodeUnitExist // such a unit already exists
	Ok
*/
func (c *Core) AddUnit(id int64) errorCodes {
	if !c.catch() {
		go c.lgr.log(errMsgCoreNotCatch).context("Unit", id).context("Method", "AddUnit").send()
		return ErrCodeCoreCatch
	}
	defer c.throw()

	if !c.storage.addUnit(id) {
		go c.lgr.log(errMsgUnitExist).context("Unit", id).context("Method", "AddUnit").send()
		return ErrCodeUnitExist
	}
	return Ok
}

/*
DelUnit - deletion of a unit.
The unit will be deleted only when all its accounts are stopped and deleted.
In case of an error, a list of not deleted accounts is returned.

Returned codes:
	ErrCodeCoreCatch // not obtained permission
	ErrCodeUnitExist // there is no such unit
	ErrCodeAccountNotStop // accounts failed to stop
	ErrCodeUnitNotEmpty // accounts are not zero
	Ok
*/
func (c *Core) DelUnit(id int64) ([]string, errorCodes) {
	if !c.catch() {
		go c.lgr.log(errMsgCoreNotCatch).context("Unit", id).context("Method", "DelUnit").send()
		return nil, ErrCodeCoreCatch
	}
	defer c.throw()

	un, ok := c.storage.delUnit(id)
	if !ok {
		go c.lgr.log(errMsgUnitNotExist).context("Unit", id).context("Method", "DelUnit").send()
		return nil, ErrCodeUnitNotExist
	}

	if accList, err := un.delAllAccounts(); err != Ok {
		go c.lgr.log(errMsgUnitNotDelAll).context("Error code", err).context("Unit", id).context("Method", "DelUnit").send()
		return accList, err
	}
	return nil, Ok
}

/*
TotalUnit - statement on all accounts of the unit.
The ID-balance array is returned.

Returned codes:
	ErrCodeCoreCatch // not obtained permission
	ErrCodeUnitExist // there is no such unit
	Ok
*/
func (c *Core) TotalUnit(id int64) (map[string]int64, errorCodes) {
	if !c.catch() {
		go c.lgr.log(errMsgCoreNotCatch).context("Unit", id).context("Method", "TotalUnit").send()
		return nil, ErrCodeCoreCatch
	}
	defer c.throw()

	un := c.storage.getUnit(id)
	if un == nil {
		go c.lgr.log(errMsgUnitNotExist).context("Unit", id).context("Method", "TotalUnit").send()
		return nil, ErrCodeUnitNotExist
	}

	return un.total(), Ok
}

/*
TotalAccount - account balance
If an account has not been created before,
it will be created with a zero balance.

Returned codes:
	ErrCodeCoreCatch // not obtained permission
	ErrCodeUnitExist // there is no such unit
	Ok
*/
func (c *Core) TotalAccount(id int64, key string) (int64, errorCodes) {
	if !c.catch() {
		go c.lgr.log(errMsgCoreNotCatch).context("Unit", id).context("Account", key).context("Method", "TotalAccount").send()
		return permitError, ErrCodeCoreCatch
	}
	defer c.throw()

	un := c.storage.getUnit(id)
	if un == nil {
		go c.lgr.log(errMsgUnitNotExist).context("Unit", id).context("Account", key).context("Method", "TotalAccount").send()
		return permitError, ErrCodeUnitNotExist
	}
	return un.getAccount(key).total(), Ok
}

/*
Start - start the application.
Only after the start you can perform transactions.
If the launch was successful, or the application is already running,
the `true` is returned, otherwise it returns `false`.
*/
func (c *Core) Start() bool {
	for i := trialLimit; i > trialStop; i-- {
		if atomic.LoadInt64(&c.hasp) == stateOpen || atomic.CompareAndSwapInt64(&c.hasp, stateClosed, stateOpen) {
			return true
		}
		runtime.Gosched()
	}
	go c.lgr.log(errMsgCoreNotStart).context("Method", "Start").send()
	return false
}

/*
Stop - stop the application.
The new transactions will not start. Old transactions are executed,
and after the end of all running transactions, the answer is returned.
*/
func (c *Core) Stop() (bool, int64) {
	for i := trialLimit; i > trialStop; i-- {
		if atomic.LoadInt64(&c.hasp) == stateClosed {
			return true, stateClosed
		} else if atomic.LoadInt64(&c.counter) == 0 &&
			atomic.CompareAndSwapInt64(&c.hasp, stateOpen, stateClosed) {
			return true, stateOpen
		}
		runtime.Gosched()
	}
	go c.lgr.log(errMsgCoreNotStop).context("Method", "Stop").send()
	return false, atomic.LoadInt64(&c.hasp)
}

/*
Load - loading data from a file.
The application stops for the duration of this operation.

Returned codes:
	ErrCodeCoreStop // unable to stop app
	ErrCodeLoadReadFile // failed to download the file
	ErrCodeLoadStrToInt64 // parsing error
	Ok
*/
func (c *Core) Load(path string) (errorCodes, map[int64]string) {
	notLoad := make(map[int64]string)
	var hasp int64
	var ok bool
	// here it is possible to change the status of `hasp` ToDo: to fix

	if ok, hasp = c.Stop(); !ok { //  hasp == stateClosed ||
		go c.lgr.log(errMsgCoreNotStop).context("Method", "Load").send()
		return ErrCodeCoreStop, notLoad
	}
	//defer c.Start()
	defer func() {
		if hasp == stateOpen {
			c.Start()
		}
	}()
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		go c.lgr.log(errMsgCoreNotReadFile).context("Path", path).context("Method", "Load").send()
		return ErrCodeLoadReadFile, notLoad
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
			go c.lgr.log(errMsgCoreParseString).context("Path", path).context("String", str).context("Method", "Load").send()
			return ErrCodeLoadStrToInt64, notLoad
		}
		balance, err := strconv.ParseInt(string(a[1]), 10, 64)
		if err != nil {
			go c.lgr.log(errMsgCoreParseString).context("Path", path).context("String", str).context("Method", "Load").send()
			return ErrCodeLoadStrToInt64, notLoad
		}
		un := c.storage.getUnit(id)
		if un == nil {
			c.storage.addUnit(id)
			un = c.storage.getUnit(id)
		}
		if _, ok := un.accounts[string(a[2])]; !ok {
			un.accounts[string(a[2])] = newAccount(balance)
		} else {
			// обработка ошибки
			notLoad[id] = string(a[2])
		}

	}
	//if hasp == stateOpen {

	//}
	//if !c.Start() { // hasp == stateClosed &&
	//	go c.lgr.New(errMsgCoreNotStart).Context("Method", "Load").Write()
	//	return ErrCodeCoreStart
	//}
	return Ok, notLoad
}

/*
Save - saving data to a file.
The application stops for the duration of this operation.

Returned codes:
	ErrCodeCoreStop // unable to stop app
	ErrCodeSaveCreateFile // could not create file
	Ok
*/
func (c *Core) Save(path string) errorCodes {
	//hasp := atomic.LoadInt64(&c.hasp)
	var hasp int64
	var ok bool
	if ok, hasp = c.Stop(); !ok {
		go c.lgr.log(errMsgCoreNotStop).context("Method", "Save").send()
		return ErrCodeCoreStop
	}
	defer func() {
		if hasp == stateOpen {
			c.Start()
		}
	}()

	var buf bytes.Buffer
	for i := uint64(0); i < storageNumber; i++ {
		for id, u := range c.storage.data[i].data {
			for key, balance := range u.totalUnsafe() {
				buf.Write([]byte(fmt.Sprintf("%d%s%d%s%s%s", id, separatorSymbol, balance, separatorSymbol, key, endLineSymbol)))
			}
		}
	}

	if ioutil.WriteFile(path, buf.Bytes(), os.FileMode(0777)) != nil {
		go c.lgr.log(errMsgCoreNotCreateFile).context("Path", path).context("Method", "Save").send()
		return ErrCodeSaveCreateFile
	}
	//if hasp == stateClosed && !c.Start() {
	//	go c.lgr.New().Context("Msg", errMsgCoreNotStart).Context("Method", "Save").Write()
	//	return ErrCodeCoreStart
	//}
	return Ok
}

/*
Begin - a new transaction is created and returned.
The application stops for the duration of this operation.
*/
func (c *Core) Begin() *Transaction {
	return newTransaction(c)
}
