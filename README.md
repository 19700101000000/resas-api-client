# RESAS-API

[RESAS-API](https://opendata.resas-portal.go.jp/)'s
client and parse sql.

## Example
```bash
./resas-api -mode=get -key=<API-KEY> -path=api/v1/prefectures -out=prefectures.json
# request: https://opendata.resas-portal.go.jp/api/v1/prefectures
# output: ./prefectures.json

./resas-api -mode=get -key=<API-KEY> -path=api/v1/cities -params="prefCode = 1" -out=cities_1.json
# request: https://opendata.resas-portal.go.jp/api/v1/cities?prefCode=1
# output: ./cities_1.json

./resas-api -mode=get_cities -key=<API-KEY> -out=json/
# requests:
#   https://opendata.resas-portal.go.jp/api/v1/cities?prefCode=1
#   ...
#   https://opendata.resas-portal.go.jp/api/v1/cities?prefCode=47
# outputs:
#   ./json/cities_1.json
#   ...
#   ./json/cities_47.json

./resas-api -mode=sql -table=prefectures -in=json/prefectures.json -cols="prefName > name, prefCode > id" -out=sql/prefectures.sql
# output: ./prefectures.sql
# sql: INSERT INTO prefectures(id,name) VALUES...

./resas-api -mode=sql -table=cities -in=json/cities_1.json -cols="cityName > name, cityCode > id, prefCode > prefecture_id" -out=sql/cities_1.sql
# output: ./cities_1.sql
# sql: INSERT INTO cities(prefecture_id,id,name) VALUES...
```

## How to...

### Show help
  ```bash
  ./resas-api -h
  ```

### Arguments
- mode
  1. get
  1. get_cities
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
    -cols="prefCode, prefName > name"
    # to
    INSERT INTO table(prefCode, name)...
    ```
  - params
    - GET method params.  
      Example...
      ```bash
      -params="prefCode = 1, cityCode = 2"
      ```

### Modes
1. get
  ```bash
  ./resas-api -mode=get -key=<API-KEY> -path=<PATH> [-out=<FILE> -params=<PARAMETERS>]
  ```
1. get_cities
   ```bash
   ./resas-api -mode=get_cities -key=<API-KEY> [-out<FILE-ARTICLE>]
   ```
1. sql
  ```bash
  ./resas-api -mode=sql -table=<TYPE> -in=<JSON> [-out=<SQL> -cols=<COLS>]
  ```
