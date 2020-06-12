package transaction

// Core
// Transaction
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

/*
newTransaction - create new Transaction.
*/
func newTransaction(c *Core) *Transaction {
	t := &Transaction{
		core: c,
		up:   make([]*request, 0, usualNumTransaction),
		reqs: make([]*request, 0, usualNumTransaction),
	}

	return t
}

/*
exeTransaction - execution of a transaction.

Returned codes:

	ErrCodeUnitNotExist // unit  not exist
	ErrCodeTransactionCatch // account not catch
	ErrCodeTransactionCredit // such a unit already exists
	Ok
*/
func (t *Transaction) exeTransaction() errorCodes {
	// catch (core)
	if !t.core.catch() {
		log(errMsgCoreNotCatch).context("Method", "exeTransaction").send()

		return ErrCodeCoreCatch
	}

	defer t.core.throw()

	// fill
	if err := t.fill(); err != Ok {
		log(errMsgTransactionNotFill).context("Method", "exeTransaction").send()

		return err
	}

	// catch (accounts)
	if err := t.catch(); err != Ok {
		log(errMsgTransactionNotCatch).context("Method", "exeTransaction").send()

		return err
	}

	// addition
	for num, i := range t.reqs {
		if res := i.account.addition(i.amount); res < 0 {
			t.rollback(num)
			t.throw(len(t.reqs))
			log(errMsgAccountCredit).context("Unit", i.id).
				context("Account", i.key).context("Amount", i.amount).
				context("Method", "exeTransaction").context("Wrong balance", res).send()

			return ErrCodeTransactionCredit
		}
	}
	// throw
	t.throw(len(t.reqs))

	return Ok
}

/*
rollback - rolled back account operations.
*/
func (t *Transaction) rollback(num int) {
	for i := 0; i < num; i++ {
		t.reqs[i].account.addition(-t.reqs[i].amount)
	}
}

/*
fill - getting accounts in the list.

Returned codes:

	ErrCodeUnitNotExist // unit not exist
	Ok
*/
func (t *Transaction) fill() errorCodes {
	for i, r := range t.reqs {
		a, err := t.core.getAccount(r.id, r.key)

		if err != Ok {
			// NOTE: log in method getAccount
			return err
		}

		t.reqs[i].account = a
	}

	return Ok
}

/*
catch - obtaining permissions from accounts.

Returned codes:

	ErrCodeTransactionCatch // account not allowed operation
	Ok
*/
func (t *Transaction) catch() errorCodes {
	for i, r := range t.reqs {
		if !r.account.catch() {
			t.throw(i)
			log(errMsgAccountNotCatch).context("Unit", r.id).
				context("Account", r.key).context("Method", "Transaction.catch").
				context("Acc counter", r.account.counter).
				context("Acc balance", r.account.balance).send()

			return ErrCodeTransactionCatch
		}
	}

	return Ok
}

/*
throw - remove permissions in accounts.
*/
func (t *Transaction) throw(num int) {
	for i, r := range t.reqs {
		if i >= num {
			break
		}

		r.account.throw()
	}
}

/*
request - single operation data
*/
type request struct {
	id      int64
	key     string
	amount  int64
	account *account
}
