# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m v2 template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/v2 .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## Структура проекта

Приведённая в этом репозитории структура проекта является рекомендуемой, но не обязательной.

Это лишь пример организации кода, который поможет вам в реализации сервиса.

При необходимости можно вносить изменения в структуру проекта, использовать любые библиотеки и предпочитаемые структурные паттерны организации кода приложения, например:
- **DDD** (Domain-Driven Design)
- **Clean Architecture**
- **Hexagonal Architecture**
- **Layered Architecture**

## Сборка

Не забыть добавить go и shortenertest в PATH:
```bash
source ~/.profile
```

Запустить сервер:
```bash
go build -o cmd/shortener/shortener cmd/shortener/*go \
    && ./cmd/shortener/shortener \
        --enable-pprof=1 \
        --audit-file=audit.log \
        --audit-url=http://localhost:8082/api/audit-logs
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

## Форматирование

```bash
go install golang.org/x/tools/cmd/goimports@latest
goimports -local "github.com/aleksandrpnshkn/go-shortener" -w ./..
```

## Документация
```bash
go install -v golang.org/x/tools/cmd/godoc@latest 

# пример доки для одного из пакетов
go doc -all ./internal/middlewares/compress/ -play

# web UI для документации
godoc -http=:6060
# По умолчанию godoc не отображает пакеты, расположенные в поддиректориях internal. 
# Чтобы увидеть служебные пакеты, добавьте в браузере параметр ?m=all: например, http://localhost:6060/pkg/?m=all.
```

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
Build ID: ad084d3320aa642fbc59d8ffcc537149cb6ceb34
Type: inuse_space
Time: 2025-11-05 00:07:21 +04
Duration: 31.11s, Total samples = 2066.37kB 
Showing nodes accounting for 2066.37kB, 100% of 2066.37kB total
      flat  flat%   sum%        cum   cum%
 1024.20kB 49.57% 49.57%  1024.20kB 49.57%  internal/profile.(*Profile).postDecode
  528.17kB 25.56% 75.13%   528.17kB 25.56%  net/http.init.func15
     514kB 24.87%   100%      514kB 24.87%  bufio.NewReaderSize (inline)
         0     0%   100%      514kB 24.87%  bufio.NewReader (inline)
         0     0%   100%  1024.20kB 49.57%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.NewAuthMiddleware.func4.1
         0     0%   100%  1024.20kB 49.57%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.NewCompressMiddleware.func3.1
         0     0%   100%  1024.20kB 49.57%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.NewDecompressMiddleware.func2.1
         0     0%   100%  1024.20kB 49.57%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.NewLogMiddleware.func1.1
         0     0%   100%  1024.20kB 49.57%  github.com/go-chi/chi/v5.(*Mux).Mount.func1
         0     0%   100%  1024.20kB 49.57%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0%   100%  1024.20kB 49.57%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0%   100%  1024.20kB 49.57%  github.com/go-chi/chi/v5/middleware.NoCache.func1
         0     0%   100%  1024.20kB 49.57%  internal/profile.Parse
         0     0%   100%  1024.20kB 49.57%  internal/profile.parseUncompressed
         0     0%   100%   528.17kB 25.56%  net/http.(*Request).write
         0     0%   100%  1538.21kB 74.44%  net/http.(*conn).serve
         0     0%   100%   528.17kB 25.56%  net/http.(*persistConn).writeLoop
         0     0%   100%   528.17kB 25.56%  net/http.(*transferWriter).doBodyCopy
         0     0%   100%   528.17kB 25.56%  net/http.(*transferWriter).writeBody
         0     0%   100%  1024.20kB 49.57%  net/http.HandlerFunc.ServeHTTP
         0     0%   100%   528.17kB 25.56%  net/http.getCopyBuf (inline)
         0     0%   100%      514kB 24.87%  net/http.newBufioReader
         0     0%   100%  1024.20kB 49.57%  net/http.serverHandler.ServeHTTP
         0     0%   100%  1024.20kB 49.57%  net/http/pprof.collectProfile
         0     0%   100%  1024.20kB 49.57%  net/http/pprof.handler.ServeHTTP
         0     0%   100%  1024.20kB 49.57%  net/http/pprof.handler.serveDeltaProfile
         0     0%   100%   528.17kB 25.56%  sync.(*Pool).Get
```
Само по себе наличие middleware в данном случае не смущает, потому что вроде как это вызвано тем, что код хендлеров обрабатывается внутри них. Но тут видно `app.Run.NewAuthMiddleware.func4.1`. Большая часть нагрузки была на чтение, и наличие тут этой миддлвари может указывать на ошибку - при редиректах эта миддлваря не должна работать, нету смысла регистрировать пользователя и сеттить куки там.

Далее в `Flame graph` (alloc space/alloc objects) видно, что приложение дважды запрашивает юзера из базы - сперва чтобы аутентифицировать, затем чтобы с ним работать в хендлере. Т.е. миддлваря складывает в контекст для хендлеров не юзера, которого уже выгрузила из БД, а его id. А хендлеры снова грузят юзера по id из БД. Можно хранить в контексте готового юзера, если он в дальнейшем он нужен. Но на самом деле при проверке токена вообще нет смысла перепроверять в БД валидный ли это user_id. Мы ведь получили его из JWT-токена, которому мы доверяем. И хендлерам одного user_id достаточно (в приложухе ща другой инфы и нету). Итого для оптимизации можно убрать походы за юзером в БД совсем, если это не сценарий с регистрацией.

Далее в `Flame graph` (alloc space/alloc objects) в вебе видно немалое потребление памяти у обработчика логов, который отправляет события во внешний сервис `audit.(*RemoteObserver).HandleEvent`. В `top` с этим вероятно связано потребление `net/http.init.func15`. Учитывая что отправка логов сделана наивно (делает для каждого события отдельный POST-запрос), это очевидное место для оптимизации. Надо сделать батчинг для логов аудита с отправкой во внешний сервис в отдельной горутине. Также как делалось удаление урлов

Так же в `Flame graph` (alloc space/alloc objects) наглядно видно сколько отдельные вызовы потребляют памяти сами (`self`). Пробежался по некоторым вызовам и:
- убрал парочку лишних переменных
- добавил резервирование памяти слайсам, которые append'ятся, и для которых итоговый размер заранее понятен

На самом деле я в оптимизации меньше всего ориентировался на `inuse objects`, но если сделать проверку из ТЗ инкремента `go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof`, то видно:
```
File: shortener
Build ID: 1e08b850897ab6f6ec7bb3cc5e44d0a79b0a3c60
Type: inuse_space
Time: 2025-11-04 23:50:37 +04
Duration: 62.31s, Total samples = 2066.37kB 
Showing nodes accounting for -1554.23kB, 75.22% of 2066.37kB total
      flat  flat%   sum%        cum   cum%
-1024.20kB 49.57% 49.57% -1024.20kB 49.57%  internal/profile.(*Profile).postDecode
 -528.17kB 25.56% 75.13%  -528.17kB 25.56%  net/http.init.func15
    -514kB 24.87%   100%     -514kB 24.87%  bufio.NewReaderSize (inline)
  512.14kB 24.78% 75.22%   512.14kB 24.78%  github.com/jackc/pgx/v5.(*Conn).getRows
         0     0% 75.22%     -514kB 24.87%  bufio.NewReader (inline)
         0     0% 75.22% -1024.20kB 49.57%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.NewAuthMiddleware.func4.1
         0     0% 75.22% -1024.20kB 49.57%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.NewCompressMiddleware.func3.1
         0     0% 75.22%   512.14kB 24.78%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.NewCompressMiddleware.func5.1
         0     0% 75.22% -1024.20kB 49.57%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.NewDecompressMiddleware.func2.1
         0     0% 75.22%   512.14kB 24.78%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.NewDecompressMiddleware.func4.1
         0     0% 75.22% -1024.20kB 49.57%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.NewLogMiddleware.func1.1
         0     0% 75.22%   512.14kB 24.78%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.NewLogMiddleware.func3.1
         0     0% 75.22%   512.14kB 24.78%  github.com/aleksandrpnshkn/go-shortener/internal/app.Run.func1.GetURLByCode.1
         0     0% 75.22%   512.14kB 24.78%  github.com/aleksandrpnshkn/go-shortener/internal/services.(*Unshortener).Unshorten
         0     0% 75.22%   512.14kB 24.78%  github.com/aleksandrpnshkn/go-shortener/internal/store/urls.(*SQLStorage).Get
         0     0% 75.22%   512.14kB 24.78%  github.com/go-chi/chi/v5.(*ChainHandler).ServeHTTP
         0     0% 75.22% -1024.20kB 49.57%  github.com/go-chi/chi/v5.(*Mux).Mount.func1
         0     0% 75.22%  -512.06kB 24.78%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 75.22%  -512.06kB 24.78%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0% 75.22% -1024.20kB 49.57%  github.com/go-chi/chi/v5/middleware.NoCache.func1
         0     0% 75.22%   512.14kB 24.78%  github.com/jackc/pgx/v5.(*Conn).Query
         0     0% 75.22%   512.14kB 24.78%  github.com/jackc/pgx/v5.(*Conn).QueryRow (inline)
         0     0% 75.22%   512.14kB 24.78%  github.com/jackc/pgx/v5/pgxpool.(*Conn).QueryRow
         0     0% 75.22%   512.14kB 24.78%  github.com/jackc/pgx/v5/pgxpool.(*Pool).QueryRow
         0     0% 75.22% -1024.20kB 49.57%  internal/profile.Parse
         0     0% 75.22% -1024.20kB 49.57%  internal/profile.parseUncompressed
         0     0% 75.22%  -528.17kB 25.56%  net/http.(*Request).write
         0     0% 75.22% -1026.06kB 49.66%  net/http.(*conn).serve
         0     0% 75.22%  -528.17kB 25.56%  net/http.(*persistConn).writeLoop
         0     0% 75.22%  -528.17kB 25.56%  net/http.(*transferWriter).doBodyCopy
         0     0% 75.22%  -528.17kB 25.56%  net/http.(*transferWriter).writeBody
         0     0% 75.22%  -512.06kB 24.78%  net/http.HandlerFunc.ServeHTTP
         0     0% 75.22%  -528.17kB 25.56%  net/http.getCopyBuf (inline)
         0     0% 75.22%     -514kB 24.87%  net/http.newBufioReader
         0     0% 75.22%  -512.06kB 24.78%  net/http.serverHandler.ServeHTTP
         0     0% 75.22% -1024.20kB 49.57%  net/http/pprof.collectProfile
         0     0% 75.22% -1024.20kB 49.57%  net/http/pprof.handler.ServeHTTP
         0     0% 75.22% -1024.20kB 49.57%  net/http/pprof.handler.serveDeltaProfile
         0     0% 75.22%  -528.17kB 25.56%  sync.(*Pool).Get
```
Как сказано в ТЗ, раз есть отрицательные значения, значит потребление как минимум в отдельных участках кода снизилось. Правда лично для меня это всё равно недостаточно наглядно пока что. Для меня более наглядный результат - скорость обработки запросов wrk: 
- 26058 requests in 0.99m
- 145016 requests in 1.06m
