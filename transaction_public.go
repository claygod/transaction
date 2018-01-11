package transactor

// Transactor
// Transaction
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

//"errors"
//"log"
//"fmt"

type Transaction struct {
	tr *Transactor
	//down []*Request
	up   []*Request
	reqs []*Request
}

/*
func (t *Transaction) Op(customer int64, account string, count int64) *Transaction {
	if count >= 0 {
		t.up = append(t.up, &Request{id: customer, key: account, amount: count})
	} else {
		t.reqs = append(t.reqs, &Request{id: customer, key: account, amount: -count})
	}
	return t
}
*/
func (t *Transaction) Debit(customer int64, account string, count int64) *Transaction {
	t.up = append(t.up, &Request{id: customer, key: account, amount: count})
	//t.reqs = append(t.reqs, &Request{id: customer, key: account, amount: count})
	return t
}

func (t *Transaction) Credit(customer int64, account string, count int64) *Transaction {
	//t.down = append(t.down, &Request{id: customer, key: account, amount: count})
	t.reqs = append(t.reqs, &Request{id: customer, key: account, amount: -count})
	return t
}

func (t *Transaction) End() errorCodes {
	//return t.tn.executeTransaction(t)
	t.reqs = append(t.reqs, t.up...)
	return t.exeTransaction()
}

/*
func Unsafe(tn *Transactor, reqs []*Request) errorCodes {
	t := &Transaction{
		tr:   tn,
		up:   make([]*Request, 0, 0),
		reqs: reqs,
	}
	return t.exeTransaction()
}
*/

type Request struct {
	id      int64
	key     string
	amount  int64
	account *Account
}
