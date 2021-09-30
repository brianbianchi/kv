# kv

> A simple, persistent key-value store.

## Usage

```console
$ go run main.go -dbpath / -port :4000
```

|   Flag    | Default value |    Description     |
| :-------: | :-----------: | :----------------: |
| `-dbpath` |  "/tmp/data"  |  Path to leveldb.  |
|  `-port`  |    ":3000"    | Port to serve API. |

## API

|   Endpoint    |              Description              |
| :-----------: | :-----------------------------------: |
|  GET `/key`   |  Returns value of the specified key.  |
|  PUT `/key`   | Creates or replaces a key-value pair. |
| DELETE `/key` |        Deletes the given key.         |

### Todo

- make distributed
