package transaction

// Core
// Logger
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	lg "log"
)

/*
logger - prints error messages to the console
*/
type logger map[string]interface{}

/*
log - return new logger
*/
func log(msg string) logger {
	return logger{"": msg} //make(logger)
}

/*
context - add context to the log
*/
func (l logger) context(k string, v interface{}) logger {
	l[k] = v

	return l
}

func (l logger) send() {
	out := ""

	for k, value := range l {
		switch v := value.(type) {
		case int, int64:
			out += fmt.Sprintf("%s: %d. ", k, v)
		case string:
			out += fmt.Sprintf("%s: %s. ", k, v)
		}
	}

	go lg.Print(out)
}
