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
go build -o cmd/shortener/shortener cmd/shortener/*go \
    && ./cmd/shortener/shortener --enable-pprof=1
```

Запустить тест:
```bash
# template
# Параметры запуска итераций разные, можно чекнуть .github/workflows/metricstest.yml
# Репозиторий - https://github.com/Yandex-Practicum/go-autotests , там инструкция как запустить
go build -o cmd/shortener/shortener cmd/shortener/*go \
    && shortenertest -test.v -test.run=^TestIteration1$ -binary-path=./cmd/shortener/shortener

go build -o cmd/shortener/shortener cmd/shortener/*go \
    && shortenertestbeta -test.v -test.run=^TestIteration15$ \
        -binary-path=cmd/shortener/shortener \
        -database-dsn="postgres://admin:qwerty@localhost:5432/shortener?sslmode=disable"
    wipedb "postgres://admin:qwerty@localhost:5432/shortener?sslmode=disable"

# Мои тесты (count для отключения кэша, помогает отлавливать flaky-тесты)
go test -count=100 ./...
```

Работа с URLом:
```bash
curl -X POST -d 'https://practicum.yandex.ru/' -i localhost:8080

curl -X POST -H "Content-Type: application/json" -d '{"url": "https://practicum.yandex.ru/"}' --compressed -i localhost:8080/api/shorten

# с куками авторизации
curl -X POST -H "Content-Type: application/json" --cookie "auth_token=TOKEN" -d '{"url": "https://example.com/"}' --compressed -i localhost:8080/api/shorten

# список урлов юзера
curl -H "Content-Type: application/json" --cookie "auth_token=TOKEN" -i localhost:8080/api/user/urls

curl -X DELETE -H "Content-Type: application/json" --cookie "auth_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjF9.ZFQYhAk2o2DDE7PMJJcYHRgb74kcYvc-oSQ9J63elnQ" -d '["mB79DTY4", "KEvfvHAz", "kekich"]' --compressed -i localhost:8080/api/user/urls

curl -X POST -H "Content-Type: application/json" -d '[{"correlation_id": "c1", "original_url": "https://practicum.yandex.ru/"}, {"correlation_id": "c2", "original_url": "https://practicum.yandex.ru/test"}]' --compressed -i localhost:8080/api/shorten/batch

curl -i localhost:8080/EwHXdJfB


# удаление с нагрузкой
docker run -it --rm --net=host alpine/bombardier --method=DELETE --header="Content-Type: application/json" --header="Cookie: auth_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjF9.ZFQYhAk2o2DDE7PMJJcYHRgb74kcYvc-oSQ9J63elnQ" --body='["mB79DTY4", "KEvfvHAz", "kekich"]' --connections=30 --rate=100 --requests=300 http://localhost:8080/api/user/urls

# создание с нагрузкой
docker run -it --rm --net=host alpine/bombardier --method=POST --body='http://example.org/' --connections=50 --rate=800 --requests=3000 http://localhost:8080/
```

Окружение:
```bash
# установить клиент для работы с БД (psql)
apt install postgresql-client

docker compose up --detach

# с хоста
psql --host 127.0.0.1 --port 5432 --username admin --password --dbname shortener
```

Для работы с миграциями установить migrate - https://github.com/golang-migrate/migrate/tree/v4.18.3/cmd/migrate . Затем в корне проекта:
```bash
cd go-shortener/cmd/shortener

~/golang-migrate/migrate create -ext sql -dir ./internal/store/migrations -seq create_example_table

~/golang-migrate/migrate -database "postgres://admin:qwerty@localhost:5432/shortener?sslmode=disable" -path ./internal/store/migrations up
~/golang-migrate/migrate -database "postgres://admin:qwerty@localhost:5432/shortener?sslmode=disable" -path ./internal/store/migrations down

# Для очистки базы
docker compose down --volumes
```

Сгенерировать моки:
```bash
# из корня проекта
./generate-mocks.bash
```

Для тестирования сервиса аудита в докере настроен mockwire:
```bash
curl -X POST -i localhost:8082/api/audit-logs
```
При успехе будет вменяемый статус и имя совпавшего stub'а в заголовке.
И можно чекнуть логи контейнера.

## Профилирование 

Профилирую два базовых сценария. Они основные для сервиса, и в них есть новый аудит, который следовало бы оптимизировать.

Для начала подготовить тестовые данные:
```bash
psql --host 127.0.0.1 --port 5432 --username admin --password --dbname shortener --file ./dev-tools/wrk/data.sql
# при успехе напишет INSERT 0 1
```

Для нагрузки использовать https://github.com/wg/wrk
```bash
wrk --threads=2 --timeout=1s --connections=4 --duration=1m --script=./dev-tools/wrk/pprof-load.lua http://localhost:8080
```

Параллельно нагрузке запустить сбор профиля:
```bash
# снять профиль памяти
curl http://localhost:8080/debug/pprof/heap?seconds=30 > profiles/base.pprof

# снять профиль CPU
curl http://localhost:8080/debug/pprof/profile?seconds=30 > profiles/base-cpu.pprof

# посмотреть в CLI
go tool pprof profiles/base.pprof
# посмотреть в браузере
go tool pprof -http=":9090" profiles/base.pprof
```

В топе вызовов различные системные функции:
```
File: shortener
Build ID: 66b3a0482a5a23002986239f94699dabf6bf3bbe
Type: inuse_space
Time: 2025-11-04 01:12:04 +04
Duration: 30.87s, Total samples = 2.02MB 
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for -2064.46kB, 100% of 2064.46kB total
Showing top 10 nodes out of 42
      flat  flat%   sum%        cum   cum%
-1024.28kB 49.62% 49.62% -1024.28kB 49.62%  github.com/jackc/pgx/v5.(*Conn).getRows
 -528.17kB 25.58% 75.20%  -528.17kB 25.58%  net/http.init.func15
 -512.01kB 24.80%   100%  -512.01kB 24.80%  internal/poll.runtime_pollSetDeadline
         0     0%   100% -1024.28kB 49.62%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.GetURLByCode.func5
         0     0%   100% -1024.28kB 49.62%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.NewAuthMiddleware.func4.1
         0     0%   100% -1024.28kB 49.62%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.NewCompressMiddleware.func3.1
         0     0%   100% -1024.28kB 49.62%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.NewDecompressMiddleware.func2.1
         0     0%   100% -1024.28kB 49.62%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.NewLogMiddleware.func1.1
         0     0%   100% -1024.28kB 49.62%  github.com/aleksandrpnshkn/go-shortener/internal/services.(*Unshortener).Unshorten
         0     0%   100% -1024.28kB 49.62%  github.com/aleksandrpnshkn/go-shortener/internal/store/urls.(*SQLStorage).Get
```
Само по себе наличие middleware в данном случае не смущает, потому что вроде как это вызвано тем, что код хендлеров обрабатывается внутри них. Но тут видно `app.Run.NewAuthMiddleware.func4.1`. Большая часть нагрузки была на чтение, и наличие тут этой миддлвари может указывать на ошибку - при редиректах эта миддлваря не должна работать, нету смысла регистрировать пользователя и сеттить куки там.

Далее в `Flame graph` в вебе видно немалое потребление памяти у обработчика логов, который отправляет события во внешний сервис `audit.(*RemoteObserver).HandleEvent`. Учитывая что он сделан наивно, делает для каждого события отдельный POST-запрос, это выглядит как очевидное место для оптимизации
