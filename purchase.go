package transaction

// Transaction
// Facade for purchase
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

type Purchase struct {
	tn      *Transaction
	buyer   int64
	seller  int64
	money   string
	product string
	count   int64
	tax     func(int64) int64
}

func newPurchase(tn *Transaction) *Purchase {
	t := func(n int64) int64 { return n }
	p := &Purchase{tn: tn, tax: t}
	return p
}

func (p *Purchase) Buyer(b int64) *Purchase {
	p.buyer = b
	return p
}

func (p *Purchase) Seller(s int64) *Purchase {
	p.seller = s
	return p
}

func (p *Purchase) Money(m string) *Purchase {
	p.money = m
	return p
}

func (p *Purchase) Product(pr string) *Purchase {
	p.product = pr
	return p
}

func (p *Purchase) Count(c int64) *Purchase {
	p.count = c
	return p
}

func (p *Purchase) Tax(t func(int64) int64) *Purchase {
	p.tax = t
	return p
}
