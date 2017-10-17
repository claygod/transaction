package transaction

// Transaction
// API
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

const countNodes int = 65536
const trialLimit int = 20000000

type Transaction struct {
	m         sync.Mutex
	counter   int64
	nodes     [countNodes]node
	customers map[int64]*Customer
}

// New - create new transaction.
func New() Transaction {
	k := Transaction{customers: make(map[int64]*Customer)}
	for i := range k.nodes {
		k.nodes[i] = newNode()
	}
	return k
}

func (t *Transaction) AddCustomer(id int64) error {
	_, ok := t.customers[id]
	if !ok {
		t.m.Lock()
		defer t.m.Unlock()
		_, ok = t.customers[id]
		if !ok {
			t.customers[id] = newCustomer()
			return nil
		}
	}
	return errors.New("This customer already exists")
}

func (t *Transaction) Transfer() *Transfer {
	return newTransfer(t)
}

func (t *Transaction) transferDo(tr *Transfer) error {
	accFrom := t.getAccount(tr.from, tr.account)
	if accFrom == nil {
		return errors.New(fmt.Sprintf("Could not find account `%s` of user `%d`", tr.account, tr.from))
	}
	accTo := t.getAccount(tr.to, tr.account)
	if accTo == nil {
		return errors.New(fmt.Sprintf("Could not find account `%s` of user `%d`", tr.account, tr.to))
	}
	if err := accFrom.reserve(tr.count); err != nil {
		return errors.New(fmt.Sprintf("Error `reserve` in  account `%s` of user `%d`: `%s`", tr.account, tr.to, err.Error()))
	}
	accFrom.give(tr.count)
	accTo.add(tr.count)
	return nil
}

func (t *Transaction) Purchase() *Purchase {
	return newPurchase(t)
}

func (t *Transaction) purchaseDo(tr *Transfer) error {

	return nil
}

func (t *Transaction) getAccount(cus int64, acc string) *Account {
	c, ok := t.customers[cus]
	if !ok {
		return nil
	}
	return c.Account(acc)
}

func (t *Transaction) Customer(id int64) *Customer {
	c, ok := t.customers[id]
	if !ok {
		return nil //errors.New("This customer does not exist")
	}
	return c
}

func (t *Transaction) DelCustomer(id int64) error {
	_, ok := t.customers[id]
	if !ok {
		return errors.New("This customer does not exist")
	}
	//
	return nil
}

func (k *Transaction) TransactionStart(n1 uint64, n2 uint64) bool {
	key1 := uint16(n1)
	key2 := uint16(n2)
	//fmt.Print("- start ", key1, " - ", key2, "BEGIN\r\n")
	counter := trialLimit
	counter2 := trialLimit

lockingStart:

	for k.nodes[key1].freeze(n1) == false {
		key1, key2 = key2, key1
		n1, n2 = n2, n1
		runtime.Gosched()
		counter--
		if counter == 0 {
			//fmt.Print("- start ", key1, " - ", key2, "FINISH-ERROR\r\n")
			return false
		}
	}
	//atomic.AddInt64(&k.counter, 1)
	if k.nodes[key2].freeze(n2) == false {
		k.nodes[key1].unfreeze(n1)
		//atomic.AddInt64(&k.counter, -1)
		key1, key2 = key2, key1
		n1, n2 = n2, n1
		runtime.Gosched()
		counter2--
		if counter2 == 0 {
			//fmt.Print("- start ", key1, " - ", key2, "FINISH-ERROR\r\n")
			return false
		}
		goto lockingStart
	}
	//atomic.AddInt64(&k.counter, 1)
	//fmt.Print("- start ", key1, " - ", key2, "FINISH-OK\r\n")
	return true
}

func (k *Transaction) TransactionEnd(n1 uint64, n2 uint64) error {
	key1 := uint16(n1)
	key2 := uint16(n2)
	//fmt.Print("- end ", key1, " - ", key2, " \r\n")

	if err := k.nodes[key1].unfreeze(n1); err != nil {
		return err
	}
	//atomic.AddInt64(&k.counter, -1)
	if err := k.nodes[key2].unfreeze(n2); err != nil {
		return err
	}
	//atomic.AddInt64(&k.counter, -1)
	return nil
}

// node - default element for queue
type node struct {
	m    sync.Mutex
	hasp int32
	arr  map[uint64]bool
}

// newNode - create new node.
func newNode() node {
	n := node{}
	n.arr = make(map[uint64]bool)
	return n
}

func (n *node) freeze(p uint64) bool {
	n.m.Lock()
	//if n.lock() == false {
	//	return false
	//}
	if _, ok := n.arr[p]; ok {
		//n.hasp = 0 // unlock
		n.m.Unlock()
		return false
	}
	n.arr[p] = true
	//n.hasp = 0 // unlock
	n.m.Unlock()
	return true
}

func (n *node) unfreeze(p uint64) error {
	n.m.Lock()
	//if n.lock() == false {
	//	return errors.New(fmt.Sprintf("Number `%d` failed to block", p))
	//}
	if _, ok := n.arr[p]; ok {

		delete(n.arr, p)
		//n.hasp = 0 // unlock
		n.m.Unlock()
		return nil
	}
	//n.hasp = 0 // unlock
	n.m.Unlock()
	return errors.New(fmt.Sprintf("Number `%d` was not blocked", p))
}

// lock - block node
func (n *node) lock() bool {
	for i := trialLimit; i > 0; i-- {
		if n.hasp == 0 && atomic.CompareAndSwapInt32(&n.hasp, 0, 1) {
			break
		}
		if i == 5 {
			return false
		}
		runtime.Gosched()
	}
	return true
}

// unlock - unblock node
func (n *node) unlock() bool {
	for i := trialLimit; i > 0; i-- {
		if n.hasp == 1 && atomic.CompareAndSwapInt32(&n.hasp, 1, 0) {
			break
		}
		if i == 5 {
			return false
		}
		//runtime.Gosched()
	}
	return true
}
