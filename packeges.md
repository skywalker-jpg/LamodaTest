# Почему я использовал именно эти пакеты?

## Список используемых пакетов
- github.com/Masterminds/squirrel
- github.com/fatih/color
- github.com/golang-migrate/migrate/v4
- github.com/google/uuid
- github.com/labstack/echo/v4

### github.com/Masterminds/squirrel
Query builder, позволяет "собирать" эффективно собирать запросы, в том числе сложные.

### github.com/fatih/color
Удобный пакет для выделения цветом текста в консоли, использовал для выведения красивых логов и подсветки уровней (Error, Warn, Info, Debug)

### github.com/golang-migrate/migrate/v4
Предоставляет надежное решение для миграции схемы базы данных на Go

### github.com/google/uuid
Простой и эффективный пакет для генерации UUID, соответствует стандарту RFC 4122

### github.com/labstack/echo/v4
В тестовом задании, где потребовалось быстро создать прототип приложения и реализовать базовый функционал, простота Echo была предпочтительной.
