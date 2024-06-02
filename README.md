# 2024_1_Los_ping-inos
Backend-репозиторий команды Los_ping-inos 

## Ссылки
* сайт    - https://jimder.ru
* swagger - http://185.241.192.216:8085/swagger/index.html

## Как генерить swagger
* в комментариях перед ручками описать доку [в таком формате](https://github.com/swaggo/swag?tab=readme-ov-file#declarative-comments-format)
* `swag init` в терминале в репозитории internal
  * если будет ругаться на кастомные структуры, можно `swag init --parseDependency  --parseInternal -g main.go`

## Как генерить proto
  * `protoc --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative *.proto`
    * в папке с прото файлом

## Как установить `protoc` 
1. https://grpc.io/docs/protoc-installation/
2. `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

## Мини дока по репозиторию
- `config`
    - в `config.yaml` настроечки: `server`, `database` и `filesPaths` для каких-то файлов на будущее
    - в `config.go` структура конфига и парсинг; парсит `viper`
- `internal`
  - `main.go` - main.
  - `delivery` - общение с внешним миром
    - `auth.go` - всё, связанное с авторизацией
    - `consts.go` - константы
    - `interfaces.go` - интерфейсы
    - `landing.go` - ручки на получение данных
    - `runserver.go` - старт сервера
  - `pkg` - обёртки и структуры для запросов - ответов
    - `requests.go` - структуры для запросов
    - `responses.go` - обёртка ответов
  - `service` - бизнес логика
    - `auth.go` - авторизация
    - `cards.go` - логика ленты
    - `interests.go` - логика интересов
    - `interfaces.go` - иНтЕрФеЙсЫ
  - `storage` - круды
    - `person.go` - person
    - `interest.go` - interest
  - `types` - кастомные гошные типы