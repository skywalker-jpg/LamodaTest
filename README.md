# Документация к тестовому заданию на позицию "Junior Go developer"

## :pushpin: Навигация

- [Техническое задание](#1)
- [Описание разработанного решения](#2)
- [Ручки](#3)
- [Как запустить локально](#4)
- [CI/CD](#5)

<a name="1"></a>

## :page_with_curl: Техническое задание
- [Тестовое задание PDF](https://drive.google.com/file/d/1Ms_jXW5iRZHg9fdQw984DBis5KQ4WYqE/view?usp=sharing)

<a name="2"></a>

## :page_facing_up: Описание разработанного решения
1. Реализована схема базы данных
    ```sql 
   CREATE TABLE IF NOT EXISTS warehouses (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        availability BOOLEAN NOT NULL
   );
   CREATE TABLE IF NOT EXISTS products (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        size VARCHAR(50) NOT NULL,
        code VARCHAR(50) UNIQUE NOT NULL
   );
   CREATE TABLE IF NOT EXISTS warehouse_product (
        id SERIAL PRIMARY KEY,
        warehouse_id INT REFERENCES warehouses(id) ON DELETE CASCADE,
        product_id INT REFERENCES products(id) ON DELETE CASCADE,
        UNIQUE (warehouse_id, product_id),
        quantity INT NOT NULL DEFAULT 0,
        reserved_quantity INT NOT NULL DEFAULT 0
   );
   ```
2. База данных заполнена тестовыми данными
```sql
INSERT INTO warehouses (name, availability) VALUES
                                                ('Warehouse A', true),
                                                ('Warehouse B', false),
                                                ('Warehouse C', true),
                                                ('Warehouse D', true);

INSERT INTO products (name, size, code) VALUES
                                            ('Product 1', 'Small', '123'),
                                            ('Product 2', 'Medium', '456'),
                                            ('Product 3', 'Large', '789'),
                                            ('Product 4', 'Small', '444');

INSERT INTO warehouse_product (warehouse_id, product_id, quantity, reserved_quantity) VALUES
                                                                                          (1, 1, 100, 20),
                                                                                          (1, 2, 50, 10),
                                                                                          (2, 3, 75, 15);
```
3. Реализованы CRUD методы для каждый сущности и соответствующие интерфейсы для них.
4. Структура `Storage` объединяет в себе интерфейсы CRUD методов сущностей
```go
type Storage struct {
	WarehouseStorage
	ProductStorage
	WarehouseProductStorage
}
```
5. Реализован кастомный хендлер для логгера `log/slog`, позволяющий выводить логи в читаемом формате и с цветовыми обозначениями для разных категорий (INFO, WARN, DEBUG, ERROR)
6. Функция `run()` получает конфиг, инициализирует на его основе логгер, базу данных, поднимает миграции, заполняет базу тестовыми данными (см. Тестовые данные TODO) и запускает http сервер
7. Пакет web содержит в себе создание сервера на основе `echo`, соответствующие хендлеры и middleware (см. Хендлеры)
8. В `Makefile` содержится команда `make up`, поднимающая `docker compose` с необходимыми зависимостями

<a name="3"></a>

## :open_hands: Ручки

| Операции | Метод | Описание                                   | Передаваемые данные (конкретная структура в описании каждого метода) |
|---------|---|--------------------------------------------|----------------------------------------------------------------------|
| POST /products  | GetWarehouseHandler  | Получение остатков на складе               | ID склада                                                            |
| POST /reserve | ReserveProductHandler | Резервация остатков товаров на складах     | ID продукта, количество для резервации, ID склада                    |
| POST /release | ReleaseProductHandler | Освобождение резервации товаров на складах | ID продукта, количество для освобождения, ID склада                  |
| POST /block | BlockWarehouseHandler | Блокировка склада                          | ID склада                                                            |
| POST /unblock | UnblockWarehouseHandler | Разблокировка склада                       | ID склада                                                            |

### Stocks

Передаваемые данные:
```go
type WarehouseProductsDTO struct {
    WarehouseID int `json:"warehouse_id"`
}
```

1. Успешный случай
- Запрос
```shell
   curl -X POST http://0.0.0.0:8080/api/v1/products \
   --header 'Content-Type: application/json' \
   --data '{
   "warehouse_id": 2
   }'
   ```
- Ответ
```json
[
   {"id":3,"warehouse_id":2,"product_id":3,"quantity":75,"reserved_quantity":15}
]
```
2. Неверные параметры
- Запрос (строка вместо числа)
```shell
   curl -X POST http://0.0.0.0:8080/api/v1/products \
   --header 'Content-Type: application/json' \
   --data '{
   "warehouse_id": "smth"
   }'
   ```
- Ответ
```json
{"error":"Invalid request body"}
```
- Запрос (несуществующий id)
```shell
   curl -X POST http://0.0.0.0:8080/api/v1/stocks \
   --header 'Content-Type: application/json' \
   --data '{
   "warehouse_id": 123
   }'
   ```
- Ответ
```json
{"message":"Not Found"}
```

### Reserve

Передаваемые данные:
```go
type ReserveDTO struct {
	Reservations []Reserve `json:"reservations"`
}

type Reserve struct {
	Code        string `json:"code"`
	Quantity    int    `json:"quantity"`
	WarehouseID int    `json:"warehouse_id"`
}
```

1. Успешный случай
- Запрос
```shell
   curl -X POST http://0.0.0.0:8080/api/v1/reserve \
   --header 'Content-Type: application/json' \
   --data '{
  "reservations": [
    {
      "code": "123",
      "quantity": 15,
      "warehouse_id": 1
    },
    {
      "code": "456",
      "quantity": 10,
      "warehouse_id": 1
    }
  ]
}'
   ```
- Ответ
```json
{"Reserved":"OK"}
```
2. Неверные параметры
- Запрос (количество резерва больше хранимого количества)
```shell
   curl -X POST http://0.0.0.0:8080/api/v1/reserve \
   --header 'Content-Type: application/json' \
   --data '{
  "reservations": [
    {
      "code": "123",
      "quantity": 999,
      "warehouse_id": 1
    },
    {
      "code": "456",
      "quantity": 20,
      "warehouse_id": 1
    }
  ]
}'
   ```
- Ответ
```json
{"error":"Can't reserve more than have"}
```
- Запрос (несуществующий code или id)
```shell
   curl -X POST http://0.0.0.0:8080/api/v1/reserve \
   --header 'Content-Type: application/json' \
   --data '{
  "reservations": [
    {
      "code": "qweqew",
      "quantity": 10,
      "warehouse_id": 1
    },
    {
      "code": "456",
      "quantity": 20,
      "warehouse_id": 23
    }
  ]
}'
   ```
- Ответ
```json
{"error":"Unable to get stored products in warehouse: failed to get warehouse product by product code: sql: no rows in result set"}

```
- Запрос (неверный формат данных)
```shell
   curl -X POST http://0.0.0.0:8080/api/v1/reserve \
   --header 'Content-Type: application/json' \
   --data '{
  "reservations": [
    {
      "code": "123",
      "quantity": "4234",
      "warehouse_id": 1
    },
    {
      "code": "456",
      "quantity": 20,
      "warehouse_id": 1
    }
  ]
}'
   ```
- Ответ
```json
{"error":"invalid request body"}
```

### Release

Передаваемые данные:
```go
type ReleaseDTO struct {
	Releases []Release `json:"releases"`
}

type Release struct {
	Code        string `json:"code"`
	Quantity    int    `json:"quantity"`
	WarehouseID int    `json:"warehouse_id"`
}
```

1. Успешный случай
- Запрос
```shell
   curl -X POST http://0.0.0.0:8080/api/v1/release \
   --header 'Content-Type: application/json' \
   --data '{
  "releases": [
    {
      "code": "123",
      "quantity": 15,
      "warehouse_id": 1
    },
    {
      "code": "456",
      "quantity":10,
      "warehouse_id":1
    }
  ]
}'
   ```
- Ответ
```json
{"Released":"OK"}
```

2. Неуспешные случаи
- Запрос (количество освобождения больше количестве резерва)
```shell
   curl -X POST http://0.0.0.0:8080/api/v1/release \
   --header 'Content-Type: application/json' \
   --data '{
  "releases": [
    {
      "code": "123",
      "quantity": 999,
      "warehouse_id": 1
    },
    {
      "code": "456",
      "quantity":1000,
      "warehouse_id":1
    }
  ]
}'
   ```
- Ответ
```json
{"error":"Can't release more than have"}
```
- Запрос (несуществующий code или id)
```shell
   curl -X POST http://0.0.0.0:8080/api/v1/release \
   --header 'Content-Type: application/json' \
   --data '{
  "releases": [
    {
      "code": "eqwe",
      "quantity": 20,
      "warehouse_id": 1
    },
    {
      "code": "456",
      "quantity":10,
      "warehouse_id":999
    }
  ]
}'
   ```
- Ответ
```json
{"error":"Unable to get stored products in warehouse: failed to get warehouse product by product code: sql: no rows in result set"}
```
- Запрос (неверный формат данных)
```shell
   curl -X POST http://0.0.0.0:8080/api/v1/release \
   --header 'Content-Type: application/json' \
   --data '{
  "releases": [
    {
      "code": "123",
      "quantity": "ds",
      "warehouse_id": "dasd"
    },
    {
      "code": "456",
      "quantity":10,
      "warehouse_id":1
    }
  ]
}'
   ```
- Ответ
```json
{"error":"invalid request body"}
```

### Block/Unblock

Передаваемые данные:
```go
type BlockDTO struct {
	WarehouseID int `json:"warehouse_id"`
}
```

1. Успешный случай
- Запрос
```shell
curl -X POST http://0.0.0.0:8080/api/v1/unblock \
  --header 'Content-Type: application/json' \
  --data '{
  "warehouse_id": 1
  }'
   ```
- Ответ
```json
{"Unblocked":"OK"}
```
2. Неверные параметры
- Запрос (строка вместо числа)
```shell
   curl -X POST http://0.0.0.0:8080/api/v1/unblock \
   --header 'Content-Type: application/json' \
   --data '{
   "warehouse_id": "smth"
   }'
   ```
- Ответ
```json
{"error":"Invalid request body"}
```
- Запрос (несуществующий id)
```shell
   curl -X POST http://0.0.0.0:8080/api/v1/stocks \
   --header 'Content-Type: application/json' \
   --data '{
   "warehouse_id": 123
   }'
   ```
- Ответ
```json
{"message":"Not Found"}
```

<a name="4"></a>

## :hammer: Как запустить локально

### Linux, MacOS

#### Запуск

0. Склонировать репозиторий, переключиться на нужную ветку, запуллить изменения
1. Поставить необходимые зависимости (`docker`, `makefile`)
2. Запустить контейнеры комадой
```shell
make up
```

#### Остановка
1. Уничтожить контейнер командой
```shell
make down
```

### Windows(Тут вопросы, у самого ubuntu)

#### Установка окружения (docker, mysql, make)

1. Для установки `docker` необходимо:
    1. Кликнуть сюда (сразу начнется
       загрузка): [Загрузка docker](https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe)
    2. Запустить скачанный файл и установить приложение с предлагаемыми настройками
    3. Запустить приложение и дождаться, пока `docker` запустится
    4. После того как запуск завершился, можете свернуть или закрыть графическое приложение
2. Для установки `make` необходимо:
    1. Установить `chocolatey`, написав команду в терминале:
       ```shell
       Set-ExecutionPolicy Bypass -Scope Process -Force; iwr https://community.chocolatey.org/install.ps1 -UseBasicParsing | iex
       ```
    2. Установить `make`:
       ```shell
       choco install make
       ```
#### Запуск

Не отличается от других ОС

<a name="5"></a>

Проверка вручную:
1. Сборка: `go build ./...`
