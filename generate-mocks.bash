#!/usr/bin/env bash

# Прерывать скрипт при ошибке
set -e

# Прерывать если не передана переменная
set -u

CURRENT_DIR_BASENAME=`basename $PWD`

# Скрипт нужно запускать из корня проекта, чтобы не возиться с путями в командах.
# Предполагаю что репа названа как на гитхабе. Учитывать ренейминг лень.
if [[ $CURRENT_DIR_BASENAME != "go-shortener" ]]; then
  echo "Run only from root project dir"
  exit 1
fi

mockgen -destination=internal/mocks/mock_store_urls_storage.go -package=mocks -mock_names Storage=MockURLsStorage ./internal/store/urls Storage
mockgen -destination=internal/mocks/mock_store_users_storage.go -package=mocks -mock_names Storage=MockUsersStorage ./internal/store/users Storage

mockgen -destination=internal/mocks/mock_services_auther.go -package=mocks ./internal/services Auther

mockgen -destination=internal/mocks/mock_services_code_generator.go -package=mocks ./internal/services CodeGenerator

mockgen -destination=internal/mocks/mock_services_codes_reserver.go -package=mocks ./internal/services CodesReserver

echo "Finish"
