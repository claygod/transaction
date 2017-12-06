# transaction

Библиотека оперирует только целыми числами. Если вы хотите работать, к примеру, с сотыми долями (центы в долларах), умножайте всё на 100. Например, полтора доллара, это будет 150.
Ограничение на максимальный размер счёта: 2 в 63 степени (9,223,372,036,854,775,807). Для примера: на счёте может лежать сумма не больше $92,233,720,368,547,758.07

Библиотека работает в параллельном режиме и может обрабатывать миллионы запросов в секунду.
Библиотека следит за тем, чтобы параллельные запросы к одному и тому же счёту не привели к ошибочному изменению баланса этого счёта.
Дебетовые и кредитные операции со счётом можно проводить ТОЛЬКО в составе транзакции.

В библиотеке две основных сущности: юнит и аккаунт.

### Unit

- Юнитом может быть клиент, фирма и т.д.
- У юнита может быть много счетов.
- Юнит не может быть удалён, если хотя бы один из его счетов не нулевой.
- Если юниту поступает некоторая сумма на несуществующий аккаунт, такой аккаунт будет создан.

### Account

- Аккаунт служит для учёта денег, акций и т.п.
- Аккаунт обязательно принадлежит какому-либо юниту.
- Аккаунт принадлежит только одному юниту.
- На одном аккаунте ведётся только один баланс.
- Баланс исчисляется только в целых числах.

## Usage

### Создание, загрузка, сохранение

### Создание/удаление счёта

### Credit/debit of an account

Credit and debit operations with the account:

```go
t.Begin().Credit(id, "USD", 1).End()
```	

```go
t.Begin().Debit(id, "USD", 1).End()
```

### Transfer

Пример перевода одного доллара с одного счёта на другой.

```go
t.Begin().
	Credit(idFrom, "USD", 1).
	Debit(idTo, "USD", 1).
	End()
```

### Покупка

Покупка, это по сути, два одновременных перевода.
```go
// Example of buying two shares of "Apple" for $10
tr.Begin().
	Credit(buyerId, "USD", 10).Debit(sellerId, "USD", 10).
	Debit(sellerId, "APPLE", 2).Credit(buyerId, "APPLE", 2).
	End()
```

## Transactor

- New
- load ("path")
- start (counter)
- ...
- stop (counter)
- save ("path")

## F.A.Q.

Почему нет возможности провести дебетовую или кредитную операции вне транзакции, ведь это наверняка было бы быстрее?
- Мы не хотим, чтобы у пользователя было искушение отдельно проводить связанные операции и самостоятельно делать транзакцию. Кроме того, мы считаем, что в мире финансов одиночные операции, это скорее исключение, чем правило.

Зависит ли производительность вашей библиотеки от количества ядер процессора?
- Зависит от процессора (размер кэша, количество ядер, частота, поколение), от оперативной памяти (размера и скорость), от количества аккаунтов, от вида диска(HDD/SSD) при сохранении и загрузке.

У меня одноядерный процессор, стоит ли использовать в этом случае вашу библиотеку?
- Производительность библиотеки изначально очень высокая, поэтому узким горлышком в вашем приложении она точно не будет. Однако системный блок всё-таки лучше модернизировать ;-)

## ToDo

- Демо с роутером: сервер с REST (авторизация и права доступа опущены)

## Bench

- BenchmarkCreditSequence-4     	 5000000	       358 ns/op
- BenchmarkCreditParallel-4     	10000000	       230 ns/op
- BenchmarkDebitSequence-4      	 5000000	       350 ns/op
- BenchmarkDebitParallel-4      	10000000	       228 ns/op
- BenchmarkTransferSequence-4   	 3000000	       547 ns/op
- BenchmarkTransferParallel-4   	 5000000	       369 ns/op
- BenchmarkBuySequence-4        	 2000000	       935 ns/op
- BenchmarkBuyParallel-4        	 2000000	       613 ns/op

Map:

- BenchmarkDebitSingle-4        	 3000000	       647 ns/op
- BenchmarkDebitParallel-4      	 3000000	       454 ns/op
- BenchmarkCreditSingle-4       	 2000000	       784 ns/op
- BenchmarkCreditParallel-4     	 3000000	       435 ns/op
- BenchmarkTransferSingle-4     	 2000000	       929 ns/op
- BenchmarkTransferParallel-4   	 3000000	       584 ns/op

sync.Map:

- BenchmarkCreditSingle-4       	 2000000	       703 ns/op
- BenchmarkCreditParallel-4     	 3000000	       489 ns/op
- BenchmarkDebitSingle-4        	 2000000	       867 ns/op
- BenchmarkDebitParallel-4      	 3000000	       415 ns/op
- BenchmarkTransferSingle-4     	 1000000	      1073 ns/op
- BenchmarkTransferParallel-4   	 2000000	       722 ns/op
- BenchmarkBuySingle-4          	 1000000	      1867 ns/op
- BenchmarkBuyParallel-4        	 1000000	      1431 ns/op

Account:

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
