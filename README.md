# oracle_price_api-



## Clone the project

```
$ git clone https://go.googlesource.com/example
$ cd example
```

## ENV File

```
$ mv .env.exmaple .env
$ vi .env
```
Fill in the env variable. Params example attract below

```
RPC=https://mainnet.infura.io/v3/3630895f60c94b159c58e16c0680b93a
API_URL=https://min-api.cryptocompare.com/data/v2
API_KEY=26f21f23e74c9d8545fea88c2d5446811f9306dbd33b83c2cae4c83ec7fa9493
```


## Docker
```
$ docker-compose up -d
```
To running with docker compose. Unit test will run automatic after built.


## Test
```
$ go test
```
To running unit test local.
## Run the project on local

```
$ go run main.go
```
The application will listen to 3000 port.


## API Endpoint
External price API source only can provide 5000 minutes histroy.

Support coin: ETH, BTC, MATIC, BNB, LINK
### Swagger
```
http://localhost:3000/swagger/index.html#/
```

### Get Last Price of ETH
```
GET /lastPrice/ETH
```
### Get Last Price of ETH with specific timestamp
```
GET /lastPrice/ETH/1666856160 //eg:time now -1800
```

### Get token average price in a time range
```
POST /rangePrice
 {
    "token": "ETH",
    "from_time_stamp" : 1666850160, // eg:Time now - 1800s
    "to_time_stamp" : 1666856580 //eg: time now
}
```






