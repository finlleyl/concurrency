# Атомарные операции 

Атомарные операции - это операции, которые выполняются за один шаг
относительно других горутин

В пакете ```sync/atomic``` реализованы атомарные операции для типов:

```int```, ```int32```, ```int64```, ```uint32```, ```uint64```, ```uintptr```, ```unsafe.Pointer```

Ниже перечислены функции для типа ```int64```. Пакет ```sync/atomic``` содержит
все остальные атомарные операции для других типов.

```AddInt64(addr *int64, delta int64)``` - увеличивает значение на указанную величину

```CompareAndSwapInt64(addr *int64, old, new int64)``` - заменяет значение, если оно равно old

```SwapInt64(addr *int64, new int64)``` - заменяет значение на new

```LoadInt64(addr *int64)``` - возвращает текущее значение

```StoreInt64(addr *int64, new int64)``` - устанавливает значение

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

func count() {
    var counter int64

    var wg sync.WaitGroup

    // горутины увеличивают значение счётчика
    for i := 0; i < 25; i++ {
        wg.Add(1)
        go func() {
            for i := 0; i < 2000; i++ {
                atomic.AddInt64(&counter, 1)
            }
            wg.Done()
        }()
    }
    wg.Wait()
    fmt.Printf("%d ", atomic.LoadInt64(&counter))
}

func main() {
    // делаем несколько попыток
    for i := 0; i < 5; i++ {
        count()
    }
} 
```