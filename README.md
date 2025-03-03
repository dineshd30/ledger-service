# Ledger Service

The Ledger Service provides the following functionality

- Record money movements (i.e.: deposits and withdrawals)
- View current balance
- View transaction history

### Running unit tests

To run unit tests on local machine execute below command

```
go test ./...
```

### Building Ledger Service locally

To build ledger service on local machine execute below command

```
  go build -o ./bin/api ./cmd/api
```

### Running Ledger Service locally

To run ledger service on local machine execute api binary as below

```
  ./bin/api
```

### Using Ledger Service

To deposit cash into ledger use below http endpoint

```
POST http://localhost:8080/ledger/304629d2-ba1f-43df-a839-26ceb869645a/transaction
Content-Type: application/json

{
  "type": "credit",
  "description": "test transaction",
  "amount": 66.33
}
```

You should see response as below

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sun, 02 Mar 2025 18:15:01 GMT
Content-Length: 163
Connection: close

{
  "data": {
    "id": "588f6ced-0410-477b-ab32-f5224bde3cdb",
    "date": 1740939301027,
    "type": "credit",
    "description": "test transaction",
    "amount": 66.33,
    "runningBalance": 166.33
  }
}
```

To withdraw cash from ledger use below http endpoint

```
POST http://localhost:8080/ledger/304629d2-ba1f-43df-a839-26ceb869645a/transaction
Content-Type: application/json

{
  "type": "debit",
  "description": "test transaction",
  "amount": 20.01
}

```

You should see response as below

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sun, 02 Mar 2025 18:15:55 GMT
Content-Length: 162
Connection: close

{
  "data": {
    "id": "4910ee7c-71ea-44c0-97e7-96e0cc8bc5e6",
    "date": 1740939355136,
    "type": "debit",
    "description": "test transaction",
    "amount": 20.01,
    "runningBalance": 146.32
  }
}

```

To view transaction history use below http endpoint

```
GET http://localhost:8080/ledger/304629d2-ba1f-43df-a839-26ceb869645a/statement
Content-Type: application/json
```

You should see response as below

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sun, 02 Mar 2025 18:17:07 GMT
Content-Length: 472
Connection: close

{
  "data": [
    {
      "id": "90c34a12-a326-4cf6-ab9b-f750a7e7261f",
      "date": 1740939181308,
      "type": "credit",
      "description": "Initial transaction",
      "amount": 100,
      "runningBalance": 100
    },
    {
      "id": "588f6ced-0410-477b-ab32-f5224bde3cdb",
      "date": 1740939301027,
      "type": "credit",
      "description": "test transaction",
      "amount": 66.33,
      "runningBalance": 166.33
    },
    {
      "id": "4910ee7c-71ea-44c0-97e7-96e0cc8bc5e6",
      "date": 1740939355136,
      "type": "debit",
      "description": "test transaction",
      "amount": 20.01,
      "runningBalance": 146.32
    }
  ]
}
```

To view last balance use below http endpoint

```
GET http://localhost:8080/ledger/304629d2-ba1f-43df-a839-26ceb869645a/balance
Content-Type: application/json
```

You should see response as below

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sun, 02 Mar 2025 18:18:43 GMT
Content-Length: 15
Connection: close

{
  "data": 146.32
}
```

### Cleaning ledger service

To clean service from local machine execute below command

```
rm -rf ./bin
```
