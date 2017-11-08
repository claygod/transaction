package transactor

// Transactor
// Transaction
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"errors"
	//"log"
	"fmt"
)

type Transaction struct {
	tn   *Transactor
	down []*Request
	up   []*Request
}

func newTransaction(tn *Transactor) *Transaction {
	t := &Transaction{
		tn:   tn,
		down: make([]*Request, 0),
		up:   make([]*Request, 0),
	}
	return t
}

func (t *Transaction) exeTransaction() error {
	if err := t.fillTransaction(); err != nil {
		return err
	}
	if err := t.catchTransaction(); err != nil {
		return err
	}
	// credit
	for num, i := range t.down {
		if res := i.account.creditAtomicFree(i.amount); res < 0 {
			err2 := t.deCredit(t.down, num)
			t.throwRequests(t.down, num)
			return errors.New(fmt.Sprintf("User `%d`, account `%s`, could not reserve `%d`. `%s`",
				t.down[num].id, t.down[num].key, i.amount, err2.Error()))
		}
	}
	// Debit
	for _, i := range t.up {
		//log.Printf("Balance: `%d`, Debt: `%d`, DEBIT: `%d`.", i.account.balance, i.account.debt, i.amount)
		i.account.debitAtomicFree(i.amount)
	}
	// throw
	t.throwTransaction()
	return nil
}

func (t *Transaction) deCredit(r []*Request, num int) error {
	for i := 0; i < num; i++ {
		r[i].account.debitAtomicFree(r[i].amount)
	}
	return nil
}

func (t *Transaction) fillTransaction() error {
	if err := t.fillRequests(t.down); err != nil {
		return err
	}
	if err := t.fillRequests(t.up); err != nil {
		return err
	}
	return nil
}

func (t *Transaction) fillRequests(requests []*Request) error {
	for i, r := range requests {
		a, err := t.tn.getAccount(r.id, r.key)
		if err != nil {
			return err
		}
		requests[i].account = a
	}
	return nil
}

func (t *Transaction) catchTransaction() error {
	if err := t.catchRequests(t.down); err != nil {
		return err
	}
	if err := t.catchRequests(t.up); err != nil {
		t.throwRequests(t.down, len(t.down))
		return err
	}
	return nil
}

func (t *Transaction) throwTransaction() {
	t.throwRequests(t.down, len(t.down))
	t.throwRequests(t.up, len(t.up))
}

func (t *Transaction) catchRequests(requests []*Request) error {
	for i, r := range requests {
		if !r.account.catch() {
			t.throwRequests(requests, i)
			return errors.New(fmt.Sprintf("Not caught account `%s` of user `%d`", r.key, r.id))
		}
	}
	return nil
}

func (t *Transaction) throwRequests(requests []*Request, num int) {
	for i, r := range requests {
		if i >= num {
			break
		}
		r.account.throw()
	}
}

func (t *Transaction) Debit(customer int64, account string, count int64) *Transaction {
	t.up = append(t.up, &Request{id: customer, key: account, amount: count})
	return t
}

func (t *Transaction) Credit(customer int64, account string, count int64) *Transaction {
	t.down = append(t.down, &Request{id: customer, key: account, amount: count})
	return t
}

func (t *Transaction) End() error {
	//return t.tn.executeTransaction(t)
	return t.exeTransaction()
}

type Request struct {
	id      int64
	key     string
	amount  int64
	account *Account
}
