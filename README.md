# 2024_1_Los_ping-inos
Backend-репозиторий команды Los_ping-inos

## Ссылки
* сайт    - http://185.241.192.216:8081/
* swagger - http://185.241.192.216:8085/swagger/index.html

## Как генерить swagger
* в комментариях перед ручками описать доку [в таком формате](https://github.com/swaggo/swag?tab=readme-ov-file#declarative-comments-format)
* `swag init` в терминале в репозитории internal
  * если будет ругаться на кастомные структуры, можно `swag init --parseDependency  --parseInternal -g main.go`

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