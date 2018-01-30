package transaction

// Core
// Config
// Copyright © 2017-2018 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

const trialLimitConst int = 2000000000 //29999999
const trialStop int = 64
const permitError int64 = -9223372036854775806
const usualNumTransaction = 4
const endLineSymbol string = "\n"
const separatorSymbol string = ";"

type errorCodes int

var trialLimit = trialLimitConst

// Hasp state
const (
	stateClosed int64 = iota
	stateOpen
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
	ErrCodeCoreCatch
	ErrCodeCoreStart
	ErrCodeCoreStop
	ErrCodeSaveCreateFile
	ErrCodeLoadReadFile
	ErrCodeLoadStrToInt64
)

// Error messages
const (
	errMsgUnitExist     string = `This unit already exists`
	errMsgUnitNotExist  string = `This unit already not exists`
	errMsgUnitNotDelAll string = `Could not delete all accounts`
	// errMsgAccountExist        string = `This account already exists`
	// errMsgAccountNotExist     string = `This account already not exists`
	// errMsgAccountNotEmpty     string = `Account is not empty`
	// errMsgAccountNotStop      string = `Account does not stop`
	errMsgAccountNotCatch     string = `Not caught account`
	errMsgAccountCredit       string = `Credit transaction error`
	errMsgCoreNotCatch        string = `Not caught transactor`
	errMsgTransactionNotFill  string = `Not fill transaction`
	errMsgTransactionNotCatch string = `Not caught transaction`
	// ErrMsgTransactionLessZero string = `Not caught transaction` //
	errMsgCoreNotStart      string = `Core does not start`
	errMsgCoreNotStop       string = `Core does not stop`
	errMsgCoreNotLoad       string = `Core does not load`
	errMsgCoreNotSave       string = `Core does not save`
	errMsgCoreNotReadFile   string = `Core does not read file`
	errMsgCoreNotCreateFile string = `Core does not create file`
	errMsgCoreParseString   string = `Could not parse line`
)
