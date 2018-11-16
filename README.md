# go_distributed_object_storage
Demo: Distributed object storage using AWS SQS service as middleware and AWS DynamoDB as metadata store.


## Usage:

### Prerequisites

- AWS IAM access key (id and secret)
- Two AWS SQS standard queues, named `godos-test` and `godos-test-located`
- Permissions to access the above queues, including:
    - create message
    - get queue url
    - delete message
    - receive message

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

After the two services are up and running, you can store/retrieve objects like below (in Python):

```python
import requests
import datetime
import os
import requests

# test put object
with open('./test.txt', 'w') as fp:
    ts = datetime.datetime.now().strftime('%Y-%m-%d-%H-%M-%S')
    content = 'This is a test file, created at: {}'.format(ts)
    fp.write(content)

with open('./test.txt', 'rb') as fp:
    object_name = 'obj-{}'.format(ts)
    api = 'http://{}/objects/{}'.format(addr, object_name)
    resp = requests.put(api, data=fp)


# test get object
resp = requests.get(api)
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

### Overview

GET request:

![get_request](https://raw.githubusercontent.com/Leboudin/go_distributed_object_storage/master/resources/godos.001.jpeg)