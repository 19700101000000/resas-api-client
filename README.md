# RESAS-API

[RESAS-API](https://opendata.resas-portal.go.jp/)'s
client and parse sql.

## Example
```bash
./resas-api -mode=get -key=<API-KEY>\
-path=api/v1/prefectures -out=prefectures.json
```

## How to...

### Show help
  ```bash
  ./resas-api -h
  ```

### Arguments
- mode
  1. get
  1. sql
- key
  - Your API KEY.
- path
- out
  - Output file name.
- in
  - Input file name.
- table
  1. prefectures
  2. cities
- cols
  - SQL's columns.  
    Example...
    ```bash
    # from
    -cols="prefCode, prefName as name"
    # to
    INSERT INTO table(prefCode, name)...
    ```

### Modes
- get
  ```bash
  ./resas-api -mode=get -key=<API-KEY> -path=<PATH> [-out=<FILE>]
  ```

- sql
  ```bash
  ./resas-api -mode=sql -table=<TYPE> -in=<JSON> -out=<SQL> [-cols=<COLS>]
  ```
