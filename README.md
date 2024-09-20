# gophkeeper

Plan:

1. Transport GRPC with TLS
2. Embed tls cert to client
3. Storage - PgSQL

## Work Scheme

    client_cli <-> client_as_local_daemon_with_state_db <-> server ?

## Заметки

связка такая: клиент(человек/скрипт) <-> client(+in-memory bd) <-> server <-> postgresql

1. Регистрация пользователя:

```
curl -X POST http://localhost:8080/register/user1/mypass
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjY4OTEzNjYsIlVzZXJJRCI6Nn0.82TtCsDh9rjkEh8V_x6m8kQcOHiAKZywIRKU29n9AW4
```

В процессе регистрации происходить добавление пользователя в BD и генерация jwt токена для дальнейшей работы.
`user1` - пользователь, `mypass` - пароль

2. Авторизация пользователя:

```
curl -X POST http://localhost:8080/auth/user1/mypass
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjY4OTE0OTIsIlVzZXJJRCI6Nn0.jArE2TGBa-mTpuryKDFRLZQweNUsWGcqtKRh6QtrAxc
```

После авторизации пользователя происходит выдача jwt токена и sync его записей с сервера в client:
`user1` - пользователь, `mypass` - пароль

```
2024/09/20 19:00:29 Starting http server: localhost:8080
2024/09/20 19:04:52 Synced records from server: 3
```

3. Добавление данных на сервер:

```
curl -X POST -H'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjY4OTE0OTIsIlVzZXJJRCI6Nn0.jArE2TGBa-mTpuryKDFRLZQweNUsWGcqtKRh6QtrAxc' --data-binary 'testok' http://localhost:8080/add/text/test/description
STORED
```

Добавление текстовых данных в систему, по аналогии добавляется `credentials` - логин пароль, `binary` - произвольные бинарные данные, `card` - номер кредитки.
Валидация типов происходит в специально модуле `dataparser`.
`add` - добавление, `text` - тип данных, `test` - имя айтема, `description` описание(метадата). --data-binary - передача даты.

4. Получение айтема по ключу:

```
curl -X GET http://localhost:8080/get/test 
testok
```

Получение даты по ключу уже из локального стораджа (при add запись прошла на сервер и осела в локальном сторадже клиента).
`get` - получение, `test` - имя даты.

5. Получение списка всех айтемов клиента (из локального стораджа):

```
curl -X GET http://localhost:8080/listall  
Name | Type | Description
--------------------------
test | TEXT | description
test2 | TEXT | description2
```

6. Удаление айтемов по ключу:

```
curl -X DELETE -H'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjY4OTE0OTIsIlVzZXJJRCI6Nn0.jArE2TGBa-mTpuryKDFRLZQweNUsWGcqtKRh6QtrAxc' http://localhost:8080/del/test             
REMOVED
```

`del` - удаление, `test` - имя ключа.

**TODO:**

1. Еще пишу тесты
2. TUI прикопал шаблон, но не успеваю, надо разбираться, как делать переход между моделями.
3. Постараюсь доделать retry, где это нужно.
4. Если получится с TUI, по возможно получится использовать json для локальной передачи в client.
5. Все упирается во время, сейчас отчетные ревью периоды.
...
