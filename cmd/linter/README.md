# cmd/linter

Линтер, проверяющий необработанный выход из программы:
- использование встроенной функции `panic`
- вызов функций `log.Fatal`/`os.Exit` вне функции `main` пакета `main`

```bash
go build -o cmd/linter/linter cmd/linter/main.go

# протестить
cmd/linter/linter cmd/linter/testdata/test/main.go
```
