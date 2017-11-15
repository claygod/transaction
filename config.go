package transactor

// Transactor
// Config
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

// import	"errors"

const trialLimit int = 2000000000
const trialStop int = 64
const permitError int64 = -9223372036854775806
const endLineSymbol string = "\n"
const separatorSymbol string = ";"

type errorCodes int

const (
	ErrOk errorCodes = 200
)
const (
	ErrCodeUnitExist errorCodes = 400 + iota
	ErrCodeUnitNotExist
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

// Error types
/*
const (
	ErrorTypeExist    string = `Exist`
	ErrorTypeNotExist string = `Not exist`
	ErrorTypeNotEmpty string = `Not empty`
	ErrorTypeNotStop  string = `Not stop`
)
*/
// Hasp state
const (
	stateOpen = iota
	stateClosed
)

// Error level
const (
	ErrLayerAccount = 100 * iota
	ErrLayerUnit
	ErrLayerTransaction
)

// Error level
const (
	ErrLevelError   string = `ERROR`
	ErrLevelWarning string = `WARNING`
	ErrLevelNotice  string = `NOTICE`
)

// Error level
const (
	errMsgUnitExist    string = `This unit already exists`
	errMsgUnitNotExist string = `This unit already not exists`
	// errMsgAccountExist        string = `This account already exists`
	// errMsgAccountNotExist     string = `This account already not exists`
	// errMsgAccountNotEmpty     string = `Account is not empty`
	// errMsgAccountNotStop      string = `Account does not stop`
	errMsgAccountNotCatch     string = `Not caught account`
	errMsgAccountCredit       string = `Credit transaction error`
	errMsgTransactorNotCatch  string = `Not caught transactor`
	errMsgTransactionNotFill  string = `Not fill transaction`
	errMsgTransactionNotCatch string = `Not caught transaction`
)

// Error_UnitExist := errors.New("This unit already exists")
