# Защита от SQL Injections

В нашем проекте мы используем библиотеку squirrel, которая поддерживает параметризацию запросов и помогает создавать SQL запросы из составных частей. 

Пример кода:
```go
query := stBuilder.
		Select(personFields).
		From(PersonTableName).
		Where(qb.And{whereMap, qb.Like{"LOWER(name)": filter.Name}}).
		RunWith(storage.dbReader)
```

# Работа с БД через сервисную учетную запись

Скрипт для создания сервисной учетной записи находится в файле **service_account.sql**

В конфигурации приложения параметры для подключения к СУБД исправлены на работу через сервисную учетную запись

```go
const DATABASE_URL string = "postgres://jimder_service_account:iamoutoftouch888@postgres:5432/JIMDER"
```

Скрипт на создание пользователя и прав лежит в папке db/hw2

# Пулл соединений и параметры соединений

Мы создаем пул соединений с помощью пакета go-sql-driver.

```go
db, err := sql.Open("postgres", psqInfo)
	if err != nil {
		logger.Logger.Fatalf("can't open db: %v", err.Error())
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
```

Значение max_connections в postgresql.conf должно быть чуть больше, чем максимальное количество содениений в пуле соединений. Таким образом, всегда будет несколько доступных соединений для прямого подключения для обслуживания и мониторинга системы

Параметр listen_adresses указывает TCP/IP адреса, по которым сервер должен прослушивать соединения от клиентских приложений. У нас значение данного параметра указано как `'*'`, что соответствует всем доступным IP-интерфейсам.

# Настройка параметров сервера и клиента

## Таймауты

Исходя из нашей бизнес логики никакой запрос не должен выполняться более 10 секунд, поэтому мы поставили таймаут в 10 секунд

```conf
statement_timeout = 10s             
lock_timeout = 10s 
```

## Логгирование и протоколирование медленных запросов

```logging_collector = on``` включает сбор в файлы.

```log_directory``` - параметр, определяющий директорию для сохранения логов, у нас это папка log.

```log_filename``` определяет формат имени лог файла, для удобства был выбран формат postgresql-%Y-%m-%d_%H%M%S.log.

```log_min_duration_statement``` определяет минимальную продолжительность sql запроса, который будет залогирован. Мы выбрали значение
4 мс т. к. для нашей бизнес логики любой запрос который выполняется дольше 4мс считается долгим.

```log_line_prefix``` определяет формат префикса для каждой строки лога. Для удобства был выбран формат '%m [%p] %q%u@%d '.

```conf
logging_collector = on            
log_directory = 'log'                  
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'      
log_min_duration_statement = 4    
log_line_prefix = '%m [%p] %q%u@%d ' 
```

## pg_stat_statements и auto_explain

Параметр ```pg_stat_statements``` предоставляет средства для отслеживания статистики выполнения sql запросов. 

Параметр ```compute_query_id``` определяет будут ли вычисляться id запросов для модуля ```pg_stat_statements```.
У нас этот параметр активен, т.к. он позволяет идентифицировать конкретные запросы, что упрощает анализ данных.

Максимальное количество уникальных sql запросов, которые могут быть отслежены с помощью ```pg_stat_statements```,
определяется параметром ```pg_stat_statements.max```. У нас оно равняется ```10000``` дабы ограничить обьем памяти, используемой для
хранения статистики.

Какие типы sql запросов должны быт отслежены определяется параметром ```pg_statements.track```, параметр ```all``` означает
что будут отслеживаться все типы запросов. 

```conf 
shared_preload_libraries = 'pg_stat_statements'
compute_query_id = on
pg_stat_statements.max = 10000
pg_stat_statements.track = all
``` 

```conf 
session_preload_libraries = 'auto_explain'
auto_explain.log_min_duration = '3s'
```