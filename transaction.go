package transaction

// Transaction
// API
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"errors"
	"fmt"
	"log"
	//"runtime"
	"sync"
	//"sync/atomic"
)

const countNodes int = 65536
const trialLimit int = 20000000
const trialStop int = 64

type Transaction struct {
	m         sync.Mutex
	customers map[int64]*Customer
}

// New - create new transaction.
func New() Transaction {
	k := Transaction{customers: make(map[int64]*Customer)}
	return k
}

func (t *Transaction) executeTransaction(o *Operation) error {
	//log.Print(1111111)
	downItems, upItems, err := t.requestToItems(o)
	if err != nil {
		return err
	}
	// Credit
	for num, i := range downItems {
		log.Printf("Balance: `%d`, Debt: `%d`, CREDIT: `%d`.", i.account.balance, i.account.debt, i.count)
		if err := i.account.reserve(i.count); err != nil {
			err2 := t.deReserve(downItems, num)
			t.throwItems(downItems, num)
			log.Printf("Balance: `%d`, Debt: `%d`, CREDIT____: `%d`.", i.account.balance, i.account.debt, i.count)
			return errors.New(fmt.Sprintf("User `%d`, account `%s`, could not reserve `%d`. `%s`",
				o.down[num].customer, o.down[num].account, i.count, err2))
		}
		log.Printf("Balance: `%d`, Debt: `%d`, CREDIT____: `%d`.", i.account.balance, i.account.debt, i.count)
	}
	for _, i := range downItems {
		i.account.give(i.count)
		log.Printf("Balance: `%d`, Debt: `%d`, give____: `%d`.", i.account.balance, i.account.debt, i.count)
	}
	// Debit
	for _, i := range upItems {
		log.Printf("Balance: `%d`, Debt: `%d`, DEBIT: `%d`.", i.account.balance, i.account.debt, i.count)
		i.account.topup(i.count)
		log.Printf("Balance: `%d`, Debt: `%d`, DEBIT______: `%d`.", i.account.balance, i.account.debt, i.count)
	}
	t.throwItems(downItems, len(downItems))
	t.throwItems(upItems, len(upItems))
	return nil
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
	accTo.topup(tr.count)
	return nil
}

func (t *Transaction) Purchase() *Purchase {
	return newPurchase(t)
}

func (t *Transaction) purchaseDo(p *Purchase) error {
	moneyBuyer := t.getAccount(p.buyer, p.money)
	if moneyBuyer == nil {
		return errors.New(fmt.Sprintf("Could not find account `%s` of user `%d`", p.money, p.buyer))
	}
	moneySeller := t.getAccount(p.seller, p.money)
	if moneySeller == nil {
		return errors.New(fmt.Sprintf("Could not find account `%s` of user `%d`", p.money, p.seller))
	}
	return nil
}

func (t *Transaction) getAccount(id int64, key string) *Account {
	c, ok := t.customers[id]
	if !ok {
		return nil
	}
	return c.Account(key)
}

func (t *Transaction) catchAccount(id int64, key string) *Account {
	c, ok := t.customers[id]
	if !ok {
		return nil
	}
	return c.catchAccount(key)
}

func (t *Transaction) Begin() *Operation {
	return newOperation(t)
}

func (t *Transaction) Customer(cid int64) *Customer {
	c, ok := t.customers[cid]
	if !ok {
		return nil //errors.New("This customer does not exist")
	}
	return c
}

func (t *Transaction) AccountStore(cid int64, key string) (int64, int64, error) {
	c, ok := t.customers[cid]
	if !ok {
		return -2, -2, errors.New("There is no such customer")
	}
	return c.AccountStore(key)
}

func (t *Transaction) DelCustomer(cid int64) error {
	_, ok := t.customers[cid]
	if !ok {
		return errors.New("This customer does not exist")
	}
	//
	return nil
}

func (t *Transaction) deReserve(items []*Item, num int) error {
	for i := 0; i < num; i++ {
		log.Print(num)
		items[i].account.unreserve(items[i].count)
	}
	return nil
}

func (t *Transaction) requestToItems(o *Operation) ([]*Item, []*Item, error) {
	downItems := make([]*Item, 0, len(o.down))
	for num, ch := range o.down {
		a := t.catchAccount(ch.customer, ch.account)
		if a != nil {
			downItems = append(downItems, &Item{account: a, count: ch.count})
		} else {
			t.throwItems(downItems, num)
			return nil, nil, errors.New(fmt.Sprintf("Could not find account `%s` of user `%d`", ch.account, ch.customer))
		}
	}
	upItems := make([]*Item, 0, len(o.up))
	for num, ch := range o.up {
		a := t.catchAccount(ch.customer, ch.account)
		if a != nil {
			upItems = append(upItems, &Item{account: a, count: ch.count})
		} else {
			t.throwItems(downItems, len(o.down))
			t.throwItems(upItems, num)
			return nil, nil, errors.New(fmt.Sprintf("Could not find account `%s` of user `%d`", ch.account, ch.customer))
		}
	}

	return downItems, upItems, nil
}

func (t *Transaction) throwItems(items []*Item, num int) {
	for i, item := range items {
		if i >= num {
			break
		}
		item.account.throw()
	}
}
