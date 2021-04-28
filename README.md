# urlshortener

URL shortener with no persistance storage :sweat_smile:.
Map is used as a temporary storage while api server is running.

## Requirements

- go 1.16

## Build

```go
    go build -o shortener cmd/shortener/main.go
```

## Usage

Binary runs on port 8080 by default.

```bash
./shortener
```

### Create User

Create user and login to API

```bash
curl --location --request POST 'http://localhost:8080/api/user' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "username": "test",
        "password": "helloworld!!!"
    }'
curl --location --request POST 'http://localhost:8080/api/login' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "username": "test",
        "password": "helloworld!!!"
    }'
```

Copy the token from login response

```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjIiLCJleHAiOjE2MTk1Nzg2NzB9.UxasiruFDy7oYwWTa4GUIySyLL9RLO5bsoxFJpgvuFk"
}
```

### Shorten the URL

Use the login in Authorization header with Bearer method

```bash
curl --location --request POST 'http://localhost:8080/api/url' \
    --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjIiLCJleHAiOjE2MTk1Nzg2NzB9.UxasiruFDy7oYwWTa4GUIySyLL9RLO5bsoxFJpgvuFk' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "real": "http://example.com"
    }'
```

Use the shortened URI from the response

```json
{
    "owner_id": 1,
    "shortened": "atC",
    "real": "http://example.com",
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "0001-01-01T00:00:00Z"
}
```

Request the shortened URI

```bash
curl --location --request GET 'http://localhost:8080/api/url/atC' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjIiLCJleHAiOjE2MTk1Nzg2NzB9.UxasiruFDy7oYwWTa4GUIySyLL9RLO5bsoxFJpgvuFk'
```

Response to shortened URI request

```json
{
    "owner_id": 1,
    "shortened": "atC",
    "real": "http://example.com",
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "0001-01-01T00:00:00Z"
}
```
