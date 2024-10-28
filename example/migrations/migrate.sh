#!/usr/bin/bash

cd "$(dirname "$0")" || exit 1

docker run --network host -v "$(pwd)":/migrations --rm -it migrate/migrate -verbose -source file://migrations -database "${DATABASE_URL}" $@
