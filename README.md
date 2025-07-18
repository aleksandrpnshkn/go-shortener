# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## Сборка

Не забыть добавить go и shortenertest в PATH:
```bash
source ~/.profile
```

Запустить сервер:
```bash
cd cmd/shortener

go build -o shortener *.go \
    && ./shortener
```

Запустить тест:
```bash
# template
# Параметры запуска итераций разные, можно чекнуть .github/workflows/metricstest.yml
go build -o shortener *.go \
    && shortenertest -test.v -test.run=^TestIteration1$ -binary-path=./shortener

# Мои тесты
go test ./...
```

Работа с URLом:
```bash
curl -X POST -d 'https://practicum.yandex.ru/' -i localhost:8080

curl -X POST -H "Content-Type: application/json" -d '{"url": "https://practicum.yandex.ru/"}' --compressed -i localhost:8080/api/shorten

curl -i localhost:8080/EwHXdJfB
```
