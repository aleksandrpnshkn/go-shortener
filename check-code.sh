#!/usr/bin/env bash

# Прерывать скрипт при ошибке
set -e

# Прерывать если не передана переменная
set -u

# Показывать запускаемые команды
set -x

go test -count=1 ./...

go build -o cmd/linter/linter cmd/linter/main.go

cmd/linter/linter ./cmd/shortener/...
cmd/linter/linter ./internal/...

go vet ./...

# go install golang.org/x/tools/cmd/goimports@latest
goimports -local "github.com/aleksandrpnshkn/go-shortener" -l -e ./cmd/
goimports -local "github.com/aleksandrpnshkn/go-shortener" -l -e ./internal/

echo "Finished\n"
