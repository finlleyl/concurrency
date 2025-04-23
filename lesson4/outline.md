# Паттерны многопотчности (продолжение)

# Паттерн конвейер

## Идея

* Разбить сложную работу на цепочку последовательных этапов
* Каждый этап: принимает поток данных, обрабатывает и отдаёт результаты следующему

## Шаги реализации

1. **Генератор**: порождает исходные данные в канал
2. **Этапы** (```add```, ```multiply```...)
   * Функция принимает два параметра:
     * ```doneCh chan struct{}``` - сигнал отмены.
     * ```inputCh chan int``` - канал входных данных.
   * Возвращает ```chan int``` с результатами
   * Внутри запускает горутину
        ```go
        go func() {
            defer close(outputCh)
            for v := range inputCh {
                 select {
                     case <-doneCh: return
                     case outputCh <- transform(v):
                 }
            }
        }()
        ```
3. **Сборка конвейера в main**
```go
doneCh := make(chan struct{})
defer close(doneCh)

inputCh := generator(doneCh, inputSlice)
addCh   := add(doneCh, inputCh)
outCh   := multiply(doneCh, addCh)

for r := range outCh {
  fmt.Println(r)
}

```

## Fan-out

Когда один этап становится узким местом, запускают несколько рабочих
для этого этапа

* fanOut принимает inputCh и число воркеров N
* Для каждого воркера вызывается add(doneCh, inputCh) - получается массив каналов []chan int
* Технически:
```go
func fanOut(doneCh chan struct{}, inputCh chan int) []chan int {
  workers := make([]chan int, N)
  for i := 0; i < N; i++ {
    workers[i] = add(doneCh, inputCh)
  }
  return workers
}

```

## Fan-in 

Обратно нужно собрать результаты из всех этих каналов в один 

* ```fanIn``` принимает ```doneCh``` и срез каналов ```...chan int```
* Для каждого входного канала запускает горутину, которая копирует данные в общий ```finalCh```
, затем сигнализирует ```WaitGroup```
* После ```wg.Wait()``` - закрывает ```finalCh```

```go
func fanIn(doneCh chan struct{}, chans ...chan int) chan int {
  var wg sync.WaitGroup
  finalCh := make(chan int)

  output := func(c chan int) {
    defer wg.Done()
    for v := range c {
      select {
      case <-doneCh: return
      case finalCh <- v:
      }
    }
  }

  wg.Add(len(chans))
  for _, c := range chans {
    go output(c)
  }

  go func() {
    wg.Wait()
    close(finalCh)
  }()

  return finalCh
}

```

## Паттерн Семафор 

### Задача 

Ограничить число горутин, одновременно работающих с ограниченным ресурсом

```go
type Semaphore struct {
ch chan struct{}
}

func NewSemaphore(max int) *Semaphore {
return &Semaphore{ch: make(chan struct{}, max)}
}

func (s *Semaphore) Acquire() { s.ch <- struct{}{} }
func (s *Semaphore) Release() { <-s.ch }

```

## Паттерн Worker Pool

### Идея

Создать фиксированное число воркеров, которые берут задачи
из общего канала очереди.

### Компоненты

* ```jobs := make(chan Job, N)``` - очередь заданий
* ```results := make(chan Result, N)``` - канал результатов
* Функция-воркер:
```go
func worker(id int, jobs <-chan int, results chan<- int) {
  for j := range jobs {
    // обработка
    results <- process(j)
  }
}

```
* Инициализация пула
```go
for w := 1; w <= numWorkers; w++ {
  go worker(w, jobs, results)
}

```
* Отправка задач в ```jobs```, затем ```close(jobs)```
* Чтение ```len(tasks)``` результатов из ```results```
