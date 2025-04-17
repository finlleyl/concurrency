# Паттерны многопоточности

## Канал всегда должен быть закрыт отправителем

Каналы всегда нужно закрывать, чтобы избежать утечки памяти или утечки
горутины. Утечка памяти происходит, когда горутина не завершается и зависает
в фоновом режине в течение всего времени работы программы 

Чтение из закрытого канала дает два параметра. Первый - это значение
```nil```. Второй - логический параметр, который сообщает, открыт
канал или нет (```true```, ```false```).

```go
package main

import "fmt"

func main() {
    ch := make(chan int)
    close(ch)

    // чтение из канала: 0, false
    read, open := <-ch

    // отправка в канал: panic
    ch <- 1

    // эта строка никогда не выполнится, так как программа будет паниковать
    fmt.Printf("Прочитанные данные: %d, Канал открыт? %t", read, open)
} 
```

Отправка данных в закрытый канал приведет к панике. Поэтому канал должен
закрывать отправитель, а не получатель

```go
package main

import "fmt"

func main() {
    // создаём канал
    ch := make(chan int)

    // вызываем горутину отправителя
    go sender(ch)

    // вызываем получателя
    recipient(ch)
}

// sender отправляет в канал числа от 0 до 9
func sender(ch chan int) {
    for i := 0; i < 10; i++ {
        ch <- i
    }

    // закрываем канал после отправки
    close(ch)
}

// recipient забирает из канала значения и выводит на экран, 
// когда канал закрыт, выходит из цикла и завершает функцию
func recipient(ch chan int) {
    // читаем данные из канала, пока он открыт
    for data := range ch {
        // и выводим их на экран
        fmt.Println(data)
    }
} 
```

## Паттерн Генератор

Паттерн Генератор генерирует данные в отдельной горутине, что позволяет
параллельно обрабатывать их и создавать новые

Это работает так: отправка и получение блокируются до тех пор,
пока отправитель и получатель не будут готовы

```go
import "fmt"

func main() {
    // данные в слайсе, которые будем отправлять
    input := []int{1, 2, 3, 4, 5, 6}

    // получаем канал с данными из генератора
    inputCh := generator(input)

    // отправляем данные потребителю через канал inputCh
    consumer(inputCh)
}

// generator — генератор, который создает канал и сразу возвращает его
func generator(input []int) chan int {
    inputCh := make(chan int)

    // через отдельную горутину генератор отправляет данные в канал
    go func() {
        // закрываем канал по завершению горутины — это отправитель
        defer close(inputCh)

        // перебираем данные в слайсе
        for _, data := range input {
            // отправляем данные в канал inputCh
            inputCh <- data
        }
    }()

    // возвращаем канал inputCh
    return inputCh
}

// consumer — потребитель проходит через канал и одновременно обрабатывает 
// данные из него (выводит на экран)
func consumer(inputCh chan int) {
    for data := range inputCh {
        fmt.Println(data)
    }
} 
```

## Паттерн обработки ошибок в горутинах 

Лучшего всего передать ошибки в основной поток и обработать их там

Создадим структуру ```Result```, которая объединит в себе
результат и ошибку, а также создадим результирующий канал ```resultCh```
и будем передавать его в основную горутину

```go
package main

import (
	"fmt"
	"errors"
)

// структура, в которую добавили ошибку
type Result struct {
	data int
	err error
}

func main() {
	// ваши данные
	input := []int{1, 2, 3, 4}

	// канал с результатами работы функции consumer
	resultCh := make(chan Result)

	// получаем канал с данными из генератора
	inputCh := generator(input)

	// порождаем горутину которая отправляет результат в resultCh вместе с ошибкой
	go consumer(inputCh, resultCh)

	// читаем результаты из канала resultCh
	for res := range resultCh {
		if res.err != nil {
			// здесь обрабатываем ошибку как обычно в Go
			log.Println("разберемся с ошибкой здесь")
		}
	}
}

// consumer вызывает другую функцию, которая возвращает ошибку
func consumer(inputCh chan int, resultCh chan Result) {
	// закрваем resultCh при завершении функции consumer
	defer close(resultCh)

	// перебираем данные из канала inputCh
	for data := range inputCh {
		// получаем ошибку
		resp, err := callDatabase(data)

		// создаем структуру
		result := Result{
			data: resp,
			err: err,
		}

		// отправляем структуру в канал
		resultCh <- result
	}
}

// generator отправляет данные в канал inputCh
func generator(input []int) chan int {
	// создаём канал, куда будем отправлять данные из слайса
	inputCh := make(chan int)

	// через отдельную горутину генератор отправляет данные в канал
	go func() {
		defer close(inputCh)

		// перебираем данные из слайса
		for _, data := range input {
			// отправляем данные в канал
			inputCh <- data
		}
	}()

	// возвращаем канал с данными
	return inputCh
}


// callDatabase просто возвращает ошибку как бы из функции обращения к базе данных
func callDatabase(data int) (int, error) {
	return data, errors.New("ошибка запроса к базе данных")
} 
```

Другой способ обработки ошибок - ```errgroup```. Если в одной 
из горутин возникает ошибка, она завершается, а ошибка возвращается.
Горутины, которые уже выполняются, продолжают работу до завершения.
Но ошибка будет возвращена только из той горутины, которая первая
её сгенерировала.

```go
package main

import (
    "context"
    "errors"
    "log"

    "golang.org/x/sync/errgroup"
)

func main() {
    // создаём переменную errgroup
    g := new(errgroup.Group)

    // наши данные
    input := []int{1, 2, 3, 4}

    // генератор возвращает канал, через который он отправляет данные
    inputCh := generator(input)

    for data := range inputCh {
        // тут объявляем новую переменную внутри цикла, чтобы копировать переменную 
        // в замыкание каждой горутины, а не использовать одно общее на всех значение.
        data := data

        // потребитель должен возвращать ошибку.
        // сигнатура анонимной функции всегда такая как в примере.
        g.Go(func() error {
            // получаем ошибку
            err := callDatabase(data)
            if err != nil {
                // возвращаем ошибку
                return err
            }

            return nil
        })
    }

    // здесь ждём выполнения горутин, и если хотя бы в одной из них возникает ошибка, 
    // то присваиваем её err и обрабатываем. В этом случае просто выводим на экран.
    // Обратите внимание, что g.Wait() ждёт завершения всех запущенных горутин, даже 
    // если приозошла ошибка.
    if err := g.Wait(); err != nil {
        log.Println(err)
    }
}

// generator возвращает канал, а затем отправляет в него данные
func generator(input []int) chan int {
    // создаём канал данных
    inputCh := make(chan int)

    // вызываем горутину в которой отправляем данные в канал inputCh
    go func() {
        // по завершении горутины закрываем канал
        defer close(inputCh)

        // перебираем данные в слайсе
        for _, data := range input {
            // отправляем данные из слайса в канал
            inputCh <- data
        }
    }()

    // возвращаем канал с данными
    return inputCh
}

// callDatabase просто возвращает ошибку
func callDatabase(data int) error {
    // допустим ошибка возникнет когда data = 3
    if data == 3 {
        return errors.New("ошибка запроса к базе данных")
    }

    return nil
} 
```

## Паттерн Стоп-Кран

Отправитель блокируется на канале до тех пор, пока получатель
не будет готов принять данные. Приём из закрытого канала
возвращает нулевое значение. Используя это свойство, можно создать
канал ```Done``` и закрыть его для завершения всех горутин

```go
package main

import (
 "log"
 "time"
)

func main() {
    input := []int{1, 2, 3, 4, 5, 6}
    handler(input)
    time.Sleep(time.Second)
}

// handler получает данные из слайса
func handler(input []int) {
    // канал для явной отмены
    doneCh := make(chan struct{})
    // когда выходим из handler — сразу закрываем канал doneCh
    defer close(doneCh)

    // теперь передаём и канал отмены doneCh
    inputCh := generator(doneCh, input)

    // забираем данные из канала
    for data := range inputCh {
        // если в данных 3 — выходим из handler
        if data == 3 {
            log.Println("Прекращаем обработку данных из канала")
            return
        }
        log.Println(data)
    }
    log.Println("Данные во входном канале закончились")
}

// generator возвращает канал с данными
func generator(doneCh chan struct{}, input []int) chan int {
    // канал, в который будем отправлять данные из слайса
    inputCh := make(chan int)

    // горутина, в которой отправляются данные в канал inputCh
    go func() {
        // по завершении закрываем канал inputCh
        defer close(inputCh)

        // перебираем данные в слайсе input
        for _, data := range input {
            select {
            // если канал doneCh закрылся - сразу выходим из горутины
            case <-doneCh:
                log.Println("Останавливаем генератор")
                return
            // отправляем данные в канал inputCh
            case inputCh <- data:
            }
        }
    }()

    // возвращаем канал с данными
    return inputCh
} 
```

Этот же паттерн можно использовать с контекстом. Используем его вместо
канала ```doneCh```

```go
package main

import (
    "log"
    "context"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    input := []int{1, 2, 3, 4, 5, 6}

    go func() {
        handler(ctx, input)
        cancel()
    }

    time.Sleep(time.Seconds)
}

// передадим контекст и данные из слайса
func handler(ctx context.Context, input []int) {
    // передаём данные и контекст в генератор
    inputCh := generator(ctx, input)

    // теперь канал для отмены не нужен

    for data := range inputCh {
        if data == 3 {
            log.Println("Прекращаем обработку данных из канала")
            return
        }
        log.Println(data)
    }
    log.Println("Данные во входном канале закончились")
}

func generator(ctx context.Context, input []int) chan int {
    inputCh := make(chan int)

    go func() {
        defer close(inputCh)

        for _, data := range input {
            select {
            // вместо отменяющего канала используем Context.Done()
            case <-ctx.Done():
                log.Println("Останавливаем генератор")
                return
            case inputCh <- data:
            }
        }
    }()

    return inputCh
} 
```