package transactor

// Transactor
// Transaction
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

//"errors"
//"log"
//"fmt"

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

func (t *Transaction) exeTransaction() errorCodes {
	if !t.tn.catch() {
		t.tn.lgr.New().Context("Msg", errMsgTransactorNotCatch).Context("Method", "exeTransaction").Write()
		return ErrCodeTransactorCatch
	}
	defer t.tn.throw()
	if err := t.fillTransaction(); err != Ok {
		t.tn.lgr.New().Context("Msg", errMsgTransactionNotFill).Context("Method", "exeTransaction").Write()
		return err
	}
	if err := t.catchTransaction(); err != Ok {
		t.tn.lgr.New().Context("Msg", errMsgTransactionNotCatch).Context("Method", "exeTransaction").Write()
		return err
	}
	// credit
	for num, i := range t.down {
		if res := i.account.creditAtomicFree(i.amount); res < 0 {
			t.deCredit(t.down, num)
			t.throwRequests(t.down, num)
			t.tn.lgr.New().Context("Msg", errMsgAccountCredit).Context("Unit", i.id).
				Context("Account", i.key).Context("Amount", i.amount).
				Context("Method", "exeTransaction").Context("Wrong balance", res).Write()
			return ErrCodeTransactionCredit
		}
	}
	// debit
	for _, i := range t.up {
		i.account.debitAtomicFree(i.amount)
	}
	// throw
	t.throwTransaction()
	return Ok
}

func (t *Transaction) deCredit(r []*Request, num int) {
	for i := 0; i < num; i++ {
		r[i].account.debitAtomicFree(r[i].amount)
	}
}

func (t *Transaction) fillTransaction() errorCodes {
	if err := t.fillRequests(t.down); err != Ok {
		return err
	}
	if err := t.fillRequests(t.up); err != Ok {
		return err
	}
	return Ok
}

func (t *Transaction) fillRequests(requests []*Request) errorCodes {
	for i, r := range requests {
		a, err := t.tn.getAccount(r.id, r.key)
		if err != Ok {
			// NOTE: log in method getAccount
			return err
		}
		requests[i].account = a
	}
	return Ok
}

func (t *Transaction) catchTransaction() errorCodes {
	if err := t.catchRequests(t.down); err != Ok {
		return err
	}
	if err := t.catchRequests(t.up); err != Ok {
		t.throwRequests(t.down, len(t.down))
		return err
	}
	return Ok
}

func (t *Transaction) throwTransaction() {
	t.throwRequests(t.down, len(t.down))
	t.throwRequests(t.up, len(t.up))
}

func (t *Transaction) catchRequests(requests []*Request) errorCodes {
	for i, r := range requests {
		if !r.account.catch() {
			t.throwRequests(requests, i)
			t.tn.lgr.New().Context("Msg", errMsgAccountNotCatch).Context("Unit", r.id).
				Context("Account", r.key).Context("Method", "catchRequests").Write()
			return ErrCodeTransactionCatch
		}
	}
	return Ok
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

func (t *Transaction) End() errorCodes {
	//return t.tn.executeTransaction(t)
	return t.exeTransaction()
}

type Request struct {
	id      int64
	key     string
	amount  int64
	account *Account
}
