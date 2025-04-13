Напиши три функции:

* ```Generate(out chan<- int)``` — пишет числа от 1 до 5.

* ```Square(in <-chan int, out chan<- int)``` — возводит в квадрат.

* ```Print(in <-chan int)``` — печатает.

Соедини их в пайплайн через каналы.