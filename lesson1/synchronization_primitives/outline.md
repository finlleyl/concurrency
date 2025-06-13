
# Примитивы синхронизации

Для синхронизации горутин можно использовать три примитива пакета ```sync```:
* ```sync.Mutex``` Мьютекс — это механизм, который позволяет
  выполнить критические участки кода только одной горутиной;
* ```sync.RWMutex``` Особый вид мьютекса, который позволяет одновременно
  выполняться либо произвольному количеству операций чтения,
  либо одной операции записи;
* ```sync.Cond``` Переменная условия, которая останавливает горутину до получения
  сигнала.

## Тип sync.Mutex

Для переменной типа ```sync.Mutex``` можно вызвать два метода:
* ```(m *Mutex) Lock()``` - блокирует мьютекс.
* ```(m *Mutex) Unlock()``` - разблокирует мьютекс.

## Тип sync.RWMutex

Этот тип мьютекса позволяет выполнять либо произвольное количество
операций чтения, либо одну операцию записи. При этом нельзя выполнять
две операции записи или одновременно запись и чтение

| Состояние        |Метод ```RLock()```|Метод ```Lock()```|
|------------------|---|---|
| Мьютекс свободен |Возьмет мьютекс в режим чтения|Возьмет мьютекс в режим записи|
|Одна горутина уже взяла мьютекс через ```RLock()```| Возьмет мьютекс в режим чтения. Блокировки не будет| Заблокирует горутину до тех пор, пока другие горутины не отпустят мьютекс|
|Одна горутина уже взяла мьютекс через ```RLock()```, вторая ждёт ```Lock()```|Заблокирует горутину до тех пор, пока другая горутина наконец не дождется метода ```Lock()``` и не отпустит его|Заблокирует горутину до тех пор, пока не будет отпущены все остальные блокировки|
|Одна горутина взяла мьютекс через ```Lock()```|Заблокирует горутину до тех пор, пока не будут отпущены все блокировки на запись|Заблокирует горутину, пока не будут отпущены все остальные блокировки|


## Тип sync.Cond

Одна горутина перменной условия блокирует сама себя, а другим горутинам нужно её
"пробудить", то есть освободить.
Переменная типа ```sync.Cond``` содержит локер-поле ```L``` типа ```sync.Locker```,
значениями которого выступают типы ```*sync.Mutex``` или ```*sync.RWMutex```. Значение
```L``` передается в функции ```sync.NewCond(l Locker) *Cond```

* ```(*Cond) Wait()``` - разблокирует локер L и вводит горутину в режим ожидания
  до получения сигнала
* ```(*Cond) Signal()``` - разблокирует одну из ожидающих горутин, если такие
  есть
* ```(*Cond) Broadcast()``` - разблокирует все горутины в очереди