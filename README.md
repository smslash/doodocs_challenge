# Doodocs Days

Техническое задание — [Doodocs](https://doodocs.kz/). Полное описание задания можете прочитать [здесь](https://doodocs.github.io/doodocs-days/backend/).

## Автор

* [@smslash](https://github.com/smslash)

## Настройка и запуск

В файле ```.env``` укажите свой mail.ru и password

```bash
SMTP_EMAIL=your_email@mail.ru
SMTP_PASSWORD=your_password
```

> Для дефолтного запуска

```console
make run
```

> Для тестов

```console
make test
```

> Через докер

```console
make build-and-run
```

## Как протестировать?

Так как реализована только серверная сторона можно отправлять запросы с помощью команды CURL

Можно также использовать другие способы:
- Postman
- HTTPie
- Браузер (RESTer, RESTClient)
- Собственные скрипты

### Роут 1

```bash
curl -X POST -F "file=@archiveName.zip" http://localhost:8080/api/archive/files
```

### Роут 2

```bash
curl -X POST -F "files[]=@cat.png" -F "files[]=@dog.png" http://localhost:8080/api/archive/information --output archive.zip
```

### Роут 3

```bash
curl -X POST -F "file=@fileName.pdf" -F "emails=example@gmail.com,example2@gmail.com" http://localhost:8080/api/mail/file
```