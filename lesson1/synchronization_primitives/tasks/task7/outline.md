# Очередь на вход

**Цель**: Есть "комната", куда одновременно может зайти только 1 горутина.

Каждая горутина ждёт, пока "освободится место", потом заходит на 200 мс и выходит. Используй Cond, чтобы сигнализировать об освобождении.