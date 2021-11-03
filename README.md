# simple device api

Simple Device API is really simple and concise to show how easy to implement a Restful Service with using Golang.

It uses [gin](https://github.com/gin-gonic/gin) framework for http router.


### Generate OpenAPI spec doc with swagger
```shell
swag init -g doc.go
```

### Build docker image locally
```shell
docker build -f build/Dockerfile -t simple-device-api:latest .
```

### Run with docker in your local
```shell
docker run -p 8080:8080 simple-device-api:latest
```

### Build and publish with docker
```shell
BUILD_VERSION=""
DOCKER_USERNAME=""
DOCKER_PASSWORD=""
sh build/build.sh
```

### Run from source code in your local
```shell
GIN_MODE=debug go run cmd/api/main.go --port 8080
```

### Swagger documentation URL
http://HOST:PORT/swagger/index.html
### Examples
#### Create a device
```shell
curl --location --request POST 'http://127.0.0.1:8080/api/v1/device' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name" : "Iphone",
    "brand": "Apple",
    "model": "13 Pro Max"
}'
```
#### Get a device
```shell
curl --location --request GET 'http://localhost:8080/api/v1/device/{id}' \
--data-raw ''
```

#### Delete a device
```shell
curl --location --request DELETE 'http://localhost:8080/api/v1/device/{id}' \
--data-raw ''
```
