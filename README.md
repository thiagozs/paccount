# Pismo Challenge

## Stack of tech

* golang 1.14.1
* SQLite3
* gin
* golangci-lin
* docker-ce

## Tools for help your life

### Httpie for curl sintact sugar

Linux based on **Debian**, you need first install **httpie**, so you can use `sudo apt install httpie`

## Project structure

```sh
├── database
│   ├── gorm.go
│   └── gorm_test.go
├── database.db
├── Dockerfile
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── LICENSE
├── main.go
├── Makefile
├── models
│   ├── account.go
│   ├── oprtype.go
│   └── transaction.go
├── README.md
└── server
    ├── controllers.go
    ├── cors.go
    ├── routes.go
    ├── server.go
    └── server_test.go
```

## Swagger API doc

Endpoint for documentation: `http://localhost:8080/swagger/index.html`

## Check lint of code

You need the `golangci-lint`, for the installation I recommend you see the documentation about that. Link for doc [Here](https://github.com/golangci/golangci-lint)

## Test project

In command line inside folder of project, run `go test ./... -v` or makefile `make test`

## Start project

Run `go run main.go` in command line.
Or you can make the binary files using `make build`, this command you going generate all binaries(***raw and zip***) files for **linux**, **macOs** and **raspBery Pie** in folder `out`.

## Generate docker image

For create a image with **Docker-CE** just running in folder of projet the command `docker build -t thiagozs/paccount .`

Runing the image after build `sudo docker run --rm --name=paccount --publish=8080:8080 thiagozs/paccount:latest`

Or just run `make image` for build the image and `make run.docker` for use the API on docker container.

## Docker Healthcheck

On the construction of image we have a change to put a little **healtcheck** on API.

```sh
Step 17/17 : HEALTHCHECK --interval=5s --timeout=2s --start-period=2s --retries=5 CMD [ "curl", "--silent", "--fail", "http://localhost:8080/ping" ]
```

## Use case

### Create account

Execute the command in your terminal `http post http://localhost:8080/accounts document_number:=123456789`

```sh
HTTP/1.1 201 Created
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE, UPDATE
Access-Control-Allow-Origin: *
Access-Control-Expose-Headers: Content-Length
Access-Control-Max-Age: 86400
Content-Length: 84
Content-Type: application/json; charset=utf-8
Date: Mon, 01 Jun 2020 05:17:05 GMT

{
    "created_at": 1590988625,
    "document_number": 123456789,
    "id": 1,
    "updated_at": 1590988625
}

```

### Find account

Execute the command in your terminal `http get http://localhost:8080/accounts/1`

```sh
HTTP/1.1 200 OK
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE, UPDATE
Access-Control-Allow-Origin: *
Access-Control-Expose-Headers: Content-Length
Access-Control-Max-Age: 86400
Content-Length: 84
Content-Type: application/json; charset=utf-8
Date: Mon, 01 Jun 2020 05:17:05 GMT

{
    "created_at": 1590988625,
    "document_number": 123456789,
    "id": 1,
    "updated_at": 1590988625
}

```

### Create transaction

Execute the command in your terminal `http post http://localhost:8080/transactions amount:=1.05 account_id:=1 operation_id:=4`

```sh
HTTP/1.1 201 Created
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE, UPDATE
Access-Control-Allow-Origin: *
Access-Control-Expose-Headers: Content-Length
Access-Control-Max-Age: 86400
Content-Length: 84
Content-Type: application/json; charset=utf-8
Date: Mon, 01 Jun 2020 05:17:05 GMT

{
    "account_id": 1,
    "amount": 1.05,
    "created_at": 1590988952,
    "id": 1,
    "operation_id": 4
}

```

### Find transaction by account

Execute the command in your terminal `http get http://localhost:8080/transactions/account/1`

```sh
HTTP/1.1 201 Created
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE, UPDATE
Access-Control-Allow-Origin: *
Access-Control-Expose-Headers: Content-Length
Access-Control-Max-Age: 86400
Content-Length: 84
Content-Type: application/json; charset=utf-8
Date: Mon, 01 Jun 2020 05:17:05 GMT

[
    {
        "account_id": 1,
        "amount": 1.05,
        "created_at": 1590988952,
        "id": 1,
        "operation_id": 4
    },
    {
        "account_id": 1,
        "amount": 3.51,
        "created_at": 1590989038,
        "id": 2,
        "operation_id": 1
    }
]

```

## Versioning and license

We use SemVer for versioning. You can see the versions available by checking the tags on this repository.

For more details about our license model, please take a look at the [LICENSE](LICENSE) file

---

2020, thiagozs
