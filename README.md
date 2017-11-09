# transaction

Варианты:
- счёт не удаляется а помечается
- у счёта счётчик, регистрирующий транзакции

Transaction:
- Начало транзакции
- Проведение операций
- Траблы с удалением аккаунтов и пользователей
- Конец транзакции

## Transactor

- New
- load ("path")
- start (counter)
- ...
- stop (counter)
- save ("path")

## Bench

- BenchmarkAccountCreditOk-4            	100000000	        15.6 ns/op
- BenchmarkAccountCreditErr-4           	100000000	        15.2 ns/op
- BenchmarkAccountCreditAtomOk-4        	100000000	        13.9 ns/op
- BenchmarkAccountCreditAtomErr-4       	100000000	        13.7 ns/op
- BenchmarkAccountCreditAtom2Ok-4       	100000000	        12.9 ns/op
- BenchmarkAccountCreditAtom2Err-4      	100000000	        12.9 ns/op
- BenchmarkAccountCreditAtomFreeOk-4    	200000000	         8.50 ns/op
- BenchmarkAccountCreditAtomFreeErr-4   	1000000000	         2.17 ns/op
- BenchmarkAccountDebitAtomFreeOk-4     	300000000	         4.64 ns/op
- BenchmarkAccountTotal-4               	2000000000	         0.26 ns/op
- BenchmarkAccountDebitOk-4             	100000000	        15.2 ns/op
- BenchmarkAccountDebitAtomOk-4         	100000000	        13.7 ns/op
- BenchmarkAccountDebitAtom2Ok-4        	100000000	        12.9 ns/op
- BenchmarkAccountAdd-4                 	100000000	        13.7 ns/op
- BenchmarkAccountAddParallel-4         	100000000	        22.8 ns/op
- BenchmarkAccountReserve-4             	100000000	        13.9 ns/op
- BenchmarkAccountReserveParallel-4     	50000000	        23.6 ns/op
- BenchmarkAccountGive-4                	100000000	        13.7 ns/op
- BenchmarkAccountGiveParallel-4        	100000000	        21.9 ns/op
- BenchmarkMapRead-4                    	100000000	        18.5 ns/op
- BenchmarkMapAdd-4                     	200000000	         9.52 ns/op
- BenchmarkSliceAdd-4                   	2000000000	         1.47 ns/op
- BenchmarkCAS-4                        	200000000	         6.42 ns/op
- BenchmarkAtomicStore-4                	300000000	         4.64 ns/op
- BenchmarkAtomicLoad-4                 	2000000000	         0.26 ns/op
- BenchmarkAtomicAdd-4                  	300000000	         4.63 ns/op
