#!/bin/bash
set -e

trap 'killall kv' SIGINT

cd $(dirname $0)

killall kv || true
sleep 0.1

go install -v
go build

./kv -dbpath=./tmp/nyc.db -addr=127.0.0.1:8080 -config=sharding.toml -shard=nyc &
./kv -dbpath=./tmp/denver.db -addr=127.0.0.1:8081 -config=sharding.toml -shard=denver &
./kv -dbpath=./tmp/sd.db -addr=127.0.0.1:8082 -config=sharding.toml -shard=sd &

wait
