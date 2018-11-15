# go_distributed_object_storage
Demo: Distributed object storage using AWS SQS service as middleware and AWS DynamoDB as metadata store.


## Usage:

### Quick start

This demo contains two parts: __API server__ and __data server__. API server provides user the ability to access objects through RESTful API,
and the data server stores objects. To make it works, you have to start them both, for example:

```sh
# shell 1
# runs API server on :8030, and listens data server on :8031
go run ./main.go  -address=:8030 -dps=:8031 server

# shell 2
# runs data server on :8031, and stores data in /var/www/godos
go run ./main.go -storage=/var/www/godos -address=:8031 dataserver
```

To see help message, you can use the following command:

`go run ./main.go -h`

help message:

```
-address string
        The server will listen on this address (default ":8030")
  -dps string
        The comma separated ip address of data provider servers, e.g. "localhost:8030,localhost:8031"
  -storage string
        The storage path will be used to store files (default "/data")
```