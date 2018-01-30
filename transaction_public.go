package transaction

// Core
// Transaction
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

/*
Transaction - preparation and execution of a transaction
*/
type Transaction struct {
	core *Core
	up   []*request
	reqs []*request
}

/*
Debit - add debit to transaction.

Input variables:
	customer - ID
	account - account string code
	count - number (type "uint64" for "less-zero" safety)
*/
func (t *Transaction) Debit(customer int64, account string, count uint64) *Transaction {
	t.up = append(t.up, &request{id: customer, key: account, amount: int64(count)})
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
	t.reqs = append(t.reqs, &request{id: customer, key: account, amount: -(int64(count))})
	return t
}

/*
End - complete the data preparation and proceed with the transaction.

Returned codes:

	ErrCodeUnitNotExist // unit  not exist
	ErrCodeTransactionCatch // account not catch
	ErrCodeTransactionCredit // such a unit already exists
	Ok
*/
func (t *Transaction) End() errorCodes {
	t.reqs = append(t.reqs, t.up...)
	return t.exeTransaction()
}
