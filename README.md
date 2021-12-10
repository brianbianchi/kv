# kv

> A simple, distributed key-value store.

## Usage

```console
$ go run -dbpath=./tmp/nyc.db -addr=127.0.0.1:8080 -config=sharding.toml -shard=nyc
```

```console
$ chmod +x ./build.sh
$ ./build.sh #example local build
```

|   Flag    |  Default value   |           Description            |
| :-------: | :--------------: | :------------------------------: |
| `-dbpath` |   "/tmp/data"    |         Path to leveldb.         |
|  `-port`  | "127.0.0.1:8080" |        Port to serve API.        |
| `-config` | "sharding.toml"  | Config file for static sharding. |
| `-shard`  | "127.0.0.1:8080" |        Name of the shard.        |

## API

|   Endpoint    |              Description              |
| :-----------: | :-----------------------------------: |
|  GET `/key`   |  Returns value of the specified key.  |
|  PUT `/key`   | Creates or replaces a key-value pair. |
| DELETE `/key` |   Deletes the given key-value pair.   |
