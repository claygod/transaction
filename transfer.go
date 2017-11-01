package transaction

// Transaction
// Facade for transfer
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

type Transfer struct {
	tn      *Transaction
	from    int64
	to      int64
	account string
	count   int64
}

func newTransfer(tn *Transaction) *Transfer {
	t := &Transfer{tn: tn}
	return t
}

func (t *Transfer) From(customer int64) *Transfer {
	t.from = customer
	return t
}

func (t *Transfer) To(customer int64) *Transfer {
	t.to = customer
	return t
}

func (t *Transfer) Account(acc string) *Transfer {
	t.account = acc
	return t
}

func (t *Transfer) Count(count int64) *Transfer {
	t.count = count
	return t
}

func (t *Transfer) Do() error {
	return t.tn.transferDo(t)
}

type Operation struct {
	tn   *Transaction
	down []Request
	up   []Request
}

func newOperation(tn *Transaction) *Operation {
	o := &Operation{
		tn:   tn,
		down: make([]Request, 0),
		up:   make([]Request, 0),
	}
	return o
}

func (o *Operation) Debit(customer int64, account string, count int64) *Operation {
	o.up = append(o.up, Request{customer, account, count})
	return o
}

func (o *Operation) Credit(customer int64, account string, count int64) *Operation {
	o.down = append(o.down, Request{customer, account, count})
	return o
}

func (o *Operation) End() error {
	return o.tn.executeTransaction(o)
}

type Request struct {
	customer int64
	account  string
	count    int64
}

type Item struct {
	account *Account
	count   int64
}
