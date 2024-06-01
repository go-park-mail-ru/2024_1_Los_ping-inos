## Для проведения нагрузочного тестирования в качестве основной сущности была выбрана сущность Person. Инструмент для нагрузочного тестрирования - wrk.

Скрипты генерят рандомные имена, емейлы и пароли.

1. Запуск нагрузочного тестирования на регистрацию пользователей: /api/v1/registration

```
wrk -t4 -c100 -d30s --latency -s db/perf_test/person_create.lua https://jimder.ru/
```

```
Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    87.40ms   48.19ms 886.07ms   88.99%
    Req/Sec   295.68     54.18   404.00     76.67%
  Latency Distribution
     50%   71.81ms
     75%   99.36ms
     90%  136.32ms
     99%  243.67ms
  35262 requests in 30.02s, 24.94MB read
Requests/sec:   1174.49
Transfer/sec:    850.64KB
```

2.  Запуск нагрузочного тестирования на получение пользователей: /api/v1/profile/{id} 

```
wrk -t4 -c100 -d30s --latency -s db/perf_test/person_get.lua https://jimder.ru/
```

```
Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    83.02ms   31.32ms 756.83ms   81.33%
    Req/Sec   303.12     58.63   470.00     77.63%
  Latency Distribution
     50%   73.60ms
     75%   85.93ms
     90%  131.87ms
     99%  195.64ms
  36203 requests in 30.03s, 38.31MB read
Requests/sec:   1205.71
Transfer/sec:      1.28MB
```
