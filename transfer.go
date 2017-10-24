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
	down []Check
	up   []Check
}

func newOperation(tn *Transaction) *Operation {
	o := &Operation{
		tn:   tn,
		down: make([]Check, 0),
		up:   make([]Check, 0),
	}
	return o
}

func (o *Operation) Debet(customer int64, account string, count int64) *Operation {
	o.up = append(o.up, Check{customer, account, count})
	return o
}

func (o *Operation) Credit(customer int64, account string, count int64) *Operation {
	o.down = append(o.down, Check{customer, account, count})
	return o
}

func (o *Operation) Do() error {
	o.tn.executeTransaction(o)
	return nil
}

type Check struct {
	customer int64
	account  string
	count    int64
}

type Item struct {
	account *Account
	count   int64
}
