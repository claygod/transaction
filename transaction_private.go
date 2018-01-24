package transactor

// Core
// Transaction
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

//"errors"
// import "log"

//"fmt"

func newTransaction(c *Core) *Transaction {
	t := &Transaction{
		core: c,
		up:   make([]*Request, 0, usualNumTransaction),
		reqs: make([]*Request, 0, usualNumTransaction),
	}
	return t
}

/**/
func (t *Transaction) exeTransaction() errorCodes {
	if !t.core.catch() {
		t.core.lgr.New().Context("Msg", errMsgCoreNotCatch).Context("Method", "exeTransaction").Write()
		return ErrCodeCoreCatch
	}
	defer t.core.throw()

	// fill
	//for i, r := range t.reqs {
	//	a, err := t.tn.getAccount(r.id, r.key)
	//	if err != Ok {
	//		t.tn.lgr.New().Context("Msg", errMsgTransactionNotFill).Context("Method", "exeTransaction").Write()
	//	}
	//	t.reqs[i].account = a
	//}
	if err := t.fill(); err != Ok {
		t.core.lgr.New().Context("Msg", errMsgTransactionNotFill).Context("Method", "exeTransaction").Write()
		return err
	}

	if err := t.catch(); err != Ok {
		t.core.lgr.New().Context("Msg", errMsgTransactionNotCatch).Context("Method", "exeTransaction").Write()
		return err
	}
	// addition
	for num, i := range t.reqs {
		if res := i.account.addition(i.amount); res < 0 {
			t.rollback(num)
			t.throw(len(t.reqs))
			t.core.lgr.New().Context("Msg", errMsgAccountCredit).Context("Unit", i.id).
				Context("Account", i.key).Context("Amount", i.amount).
				Context("Method", "exeTransaction").Context("Wrong balance", res).Write()
			return ErrCodeTransactionCredit
		}
	}
	// throw
	t.throw(len(t.reqs))
	return Ok
}

func (t *Transaction) rollback(num int) {
	for i := 0; i < num; i++ {
		//log.Print("~~~~~~~~~", i)
		t.reqs[i].account.addition(-t.reqs[i].amount)
	}
}

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

func (t *Transaction) catch() errorCodes {
	for i, r := range t.reqs {
		if !r.account.catch() {
			t.throw(i)
			t.core.lgr.New().Context("Msg", errMsgAccountNotCatch).Context("Unit", r.id).
				Context("Account", r.key).Context("Method", "catch").Write()
			return ErrCodeTransactionCatch
		}
	}
	return Ok
}

func (t *Transaction) throw(num int) {
	for i, r := range t.reqs {
		if i >= num {
			break
		}
		r.account.throw()
	}
}
