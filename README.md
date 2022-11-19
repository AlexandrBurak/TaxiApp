
# User Service for Taxi app

Auth/taxi ordering


## Installation

Clone repository

```bash
git clone github.com/zxcghoulhunter/InnoTaxi
```

Install dependencies
```bash
go mod tidy
```

Run application
```bash
go run main.go
```

OR

Use docker
```bash
docker compose up
```


## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`DB_HOST`
`DB_PORT`
`DB_USERNAME`
`DB_PASSWORD`
`DB_NAME` - Postgres for users data storage

`SECRET` - secret key for JWT auth

`DB_LOG`- Mongo for logging

`DB_CACHE` - Redis for cache

`JWT_EXPIRATION_TIME` 


