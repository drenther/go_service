# go_service

A mini golang webserver using echo

## Setup

- Install [golang](https://golang.org/dl/) from the official site
- `git clone` this repo
- cd into this repo
- then run `go get` to install all the dependencies
- then run `go build` to build the server
- now you will have a binary named `go_service` in the repo 
- run this binary and your server should be up

## API Endpoints

### POST /

Save a new task for the user

#### Headers

```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer userToken"
}
```

#### Request Body

```json
{
  "body": "Task content"
}
```

### Response Body

```json
{
  "id": "random_uuid_v4",
  "user": "userToken",
  "body": "Task Content"
}
```

### GET /

Lists all the tasks for the user

#### Headers

```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer userToken"
}
```

#### Response Body

```json
[
  {
    "id": "random_uuid_v4",
    "user": "userToken",
    "body": "Task Content"
  }
]
```

> NOTE: may return `204` HTTP STATUS CODE (no-content) if the database is empty, handle that in the frontend