# 2024_1_Los_ping-inos
Backend-репозиторий команды Los_ping-inos

## Мини дока по репозиторию
- `config`
    - в `config.json` настроечки: `server`, `database` и `filesPaths` для каких-то файлов на будущее
    - в `config.go` структура конфига и парсинг; парсит `viper`

- `internal`
  - `app`
    - `main.go` - старт сервера
  - `storage` - 50 оттенков крудов
  - `types` - кастомные гошные типы

## Добавочное
`logrus` - логгер