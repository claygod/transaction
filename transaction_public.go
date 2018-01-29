package transactor

// Core
// Transaction
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

/*
Transaction - preparation and execution of a transaction
*/
type Transaction struct {
	core *Core
	up   []*Request
	reqs []*Request
}

/*
Debit - add debit to transaction.

Input variables:
	customer - ID
	account - account string code
	count - number (type "uint64" for "less-zero" safety)
*/
func (t *Transaction) Debit(customer int64, account string, count uint64) *Transaction {
	t.up = append(t.up, &Request{id: customer, key: account, amount: int64(count)})
	return t
}

/*
Credit - add credit to transaction.

Input variables:
	customer - ID
	account - account string code
	count - number (type "uint64" for "less-zero" safety)
*/
func (t *Transaction) Credit(customer int64, account string, count uint64) *Transaction {
	t.reqs = append(t.reqs, &Request{id: customer, key: account, amount: -(int64(count))})
	return t
}

/*
End - complete the data preparation and proceed with the transaction.
*/
func (t *Transaction) End() errorCodes {
	t.reqs = append(t.reqs, t.up...)
	return t.exeTransaction()
}

/*
Request - single operation data
*/
type Request struct {
	id      int64
	key     string
	amount  int64
	account *Account
}
