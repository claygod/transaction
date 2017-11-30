package transactor

// Transactor
// Config
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

const trialLimit int = 2000000000
const trialStop int = 64
const permitError int64 = -9223372036854775806
const endLineSymbol string = "\n"
const separatorSymbol string = ";"

type errorCodes int

// Hasp state
const (
	stateOpen = iota
	stateClosed
)

// No error code
const (
	Ok errorCodes = 200
)

// Error codes
const (
	ErrCodeUnitExist errorCodes = 400 + iota
	ErrCodeUnitNotExist
	ErrCodeUnitNotEmpty
	ErrCodeAccountExist
	ErrCodeAccountNotExist
	ErrCodeAccountNotEmpty
	ErrCodeAccountNotStop
	ErrCodeTransactionFill
	ErrCodeTransactionCatch
	ErrCodeTransactionCredit
	ErrCodeTransactionDebit
	ErrCodeTransactorCatch
	ErrCodeTransactorStart
	ErrCodeTransactorStop
	ErrCodeSaveCreateFile
	ErrCodeLoadReadFile
	ErrCodeLoadStrToInt64
)

// Error messages
const (
	errMsgUnitExist    string = `This unit already exists`
	errMsgUnitNotExist string = `This unit already not exists`
	// errMsgAccountExist        string = `This account already exists`
	// errMsgAccountNotExist     string = `This account already not exists`
	// errMsgAccountNotEmpty     string = `Account is not empty`
	// errMsgAccountNotStop      string = `Account does not stop`
	errMsgAccountNotCatch         string = `Not caught account`
	errMsgAccountCredit           string = `Credit transaction error`
	errMsgTransactorNotCatch      string = `Not caught transactor`
	errMsgTransactionNotFill      string = `Not fill transaction`
	errMsgTransactionNotCatch     string = `Not caught transaction`
	errMsgTransactorNotStart      string = `Transactor does not start`
	errMsgTransactorNotStop       string = `Transactor does not stop`
	errMsgTransactorNotLoad       string = `Transactor does not load`
	errMsgTransactorNotSave       string = `Transactor does not save`
	errMsgTransactorNotReadFile   string = `Transactor does not read file`
	errMsgTransactorNotCreateFile string = `Transactor does not create file`
	errMsgTransactorParseString   string = `Could not parse line`
)
