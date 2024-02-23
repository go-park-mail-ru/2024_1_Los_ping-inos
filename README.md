# 2024_1_Los_ping-inos
Backend-репозиторий команды Los_ping-inos

## Мини дока по репозиторию
- `config`
    - в `config.json` настроечки: `server`, `database` и `filesPaths` для каких-то файлов на будущее
    - в `config.go` структура конфига и парсинг; парсит `viper`

- `internal`
  - `app`
    - `main.go` - main.
  - `delivery` - общение с внешним миром
    - `runserver.go` - старт сервера
    - `landing.go` - пример обработки страницы
  - `service` - бизнес логика
    - `getExample.go` - пример бизнес логики
  - `storage` - круды
    - `entityExample.go` - пример вытаскивания информации из бд (без бд) 
  - `types` - кастомные гошные типы

## Добавочное
`logrus` - логгер