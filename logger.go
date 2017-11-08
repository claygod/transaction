package transactor

// Transactor
// Logger
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	// "io"
	"log"
)

type logger map[string]interface{}

func (l logger) New() logger {
	x := make(logger)
	return x
}

func (l logger) Context(k string, v interface{}) logger {
	l[k] = v
	return l
}

func (l logger) Write() {
	out := ""
	for k, value := range l {
		switch v := value.(type) {
		case int, int64:
			out += fmt.Sprintf("%s: %d. ", k, v)
		case string:
			out += fmt.Sprintf("%s: %s. ", k, v)
		}
	}
	//w.Print(out)
	log.Print(out)
}
