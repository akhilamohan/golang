## Prerequisites
* go 1.19 should be installed on the setup where unit tests are executed
* make and docker should be installed on the setup to run application

## Running the application

```
cd golang
make docker-up
```
Application will be listening on port 3000

## Running unit tests

```
cd golang/pkg/controller
go test .
```

## Running the application

```
cd golang
make docker-down
```

