# Kumot Coding Challenge

## Problem Statement

Develop a golang project that has an API endpoint that takes a list of github
usernames (up to a max of 10 names) and returns to the user a list of basic
information for those users including:

<ul>
  <li>name</li>
  <li>login</li>
  <li>company</li>
  <li>number of followers</li>
  <li>number of public repos</li>
</ul>

### Objectives:

<ol>
    <li>The returned users should be sorted alphabetically by name</li>
    <li>If some usernames cannot be found, this should not fail the other usernames that
were requested</li>
    <li>Implement a caching layer that will store a user that has been retrieved from
GitHub for 2 minutes. If a user's information has been cached, it should NOT hit
Github again to retrieve the same user (if it is still in the cache window), instead it
should return the user from the cache.</li>
    <li>Include the appropriate error handling and the appropriate frameworks to make
the project more extensible.</li>
    <li>Use regular http calls to hit github's API, don’t use any pre made github Golang
libraries to integrate with github</li>
</ol>

#### Note:

The API endpoint needed to get Github user information is
`https://api.github.com/users/{username}`

## Architecture and Design

```Shell
.
├── Dockerfile                       # Dockerfile for the the golang application
├── LICENSE
├── README.md
├── cmd                              # Directory to store executables--entry point
│   └── server                       # Contains the main file which serves the backend
│       └── main.go
├── docker-compose.yml               # Simplify the deployment of the backend and redis image
├── go.mod
├── go.sum
├── pkg                              # Contains the controller, services, repository
│   ├── domain                       # Directory containing the definition of models/domain
│   │   └── user.go                  # Github user details model
│   │
│   ├── http                         # Controller Layer
│   │   └── rest
│   │       ├── controller.go        # Handler definition
│   │       ├── response.go          # Struct definition for appropriate server responses
│   │       └── user.go              # Struct definition for parsing request body
│   │
│   ├── listing                      # Service Layer: application logic of the application
│   │   ├── service.go               # Handle caching and repository calls
│   │   ├── service_test.go          # Unit and mock test for the service layer
│   │   └── sorter.go                # Define helper function for sorting the github details
│   │
│   ├── logging                      # Logging service
│   │   └── logger.go                # Uses zerolog as the logger
│   │
│   ├── repository                   # Repository Layer: Data Access and Persistence
│   │   ├── api
│   │   │   └── repository.go        # Definition of Repository interface, github implementation
│   │   └── cache
│   │       └── cache.go             # Definition of Cache Interface, Redis Implementation
│   │
│   └── router
│       └── router.go                # Definition of router interface, Chi-Router implementation
└── redis.conf                       # Basic custom configuration for our redis Cache
```

## Installation and Usage

### Setup

Make sure you have Docker installed on your device, otherwise you may [download](https://www.docker.com/products/docker-desktop/) docker from their webpage. Installing docker desktop also includes Docker CLI with Docker Compose plugin.

If you already have Docker, simply follow the steps to deploy our backend and redis.

1. Open your linux terminal and clone this repository.

```shell
$ git@github.com:jrpespinas/kumot-coding-challenge.git
```

2. Locate the project directory

```shell
$ cd ~/kumot-coding-challenge
```

3. Run docker compose

```shell
$ docker compose up
```

Docker will begin deploying your golang backend and redis cache.
You should expect the following log from your terminal after the deployment.

```shell
kumot-server  | {"level":"info","time":"2022-06-29T06:34:02Z","message":"Redis server running at redis-server:6379"}
kumot-server  | {"level":"info","time":"2022-06-29T06:34:02Z","message":"Serving at 8080"}
```

### Test

1. Open a separate linux terminal to run initial testing. Change directory to the project.

```shell
$ cd ~/kumot-coding-challenge
```

2. Go to the service layer

```shell
$ cd pkg/listing
```

3. Run test

```shell
$ go test -v
```

4. You should expect similar results below

```shell
=== RUN   TestValidateEmptyList
--- PASS: TestValidateEmptyList (0.00s)
=== RUN   TestValidateMoreThanTenUsers
--- PASS: TestValidateMoreThanTenUsers (0.00s)
=== RUN   TestShowDetailsFromGithub
{"level":"info","layer":"service","time":"2022-06-29T14:39:21+08:00","message":"checking cache for existing user details"}
{"level":"info","layer":"service","time":"2022-06-29T14:39:21+08:00","message":"pulling from user details from api"}
{"level":"info","layer":"service","time":"2022-06-29T14:39:21+08:00","message":"store user details in cache"}
{"level":"info","layer":"service","time":"2022-06-29T14:39:21+08:00","message":"appending record to slice"}
{"level":"info","layer":"service","time":"2022-06-29T14:39:21+08:00","message":"sorting user details"}
    service_test.go:108: PASS:  GetDetails()
    service_test.go:109: PASS:  Get()
    service_test.go:109: PASS:  Set()
--- PASS: TestShowDetailsFromGithub (0.00s)
=== RUN   TestShowDetailsFromRedis
{"level":"info","layer":"service","time":"2022-06-29T14:39:21+08:00","message":"checking cache for existing user details"}
{"level":"info","layer":"service","time":"2022-06-29T14:39:21+08:00","message":"appending record to slice"}
{"level":"info","layer":"service","time":"2022-06-29T14:39:21+08:00","message":"sorting user details"}
    service_test.go:154: PASS:  Get()
--- PASS: TestShowDetailsFromRedis (0.00s)
PASS
ok      github.com/jrpespinas/kumot-coding-challenge/pkg/listing        0.003s
```

### Usage

1. Open another linux terminal and paste the following curl request

```shell
curl -X POST http://localhost:8080/users -H 'Content-Type: application/json' -d '{"usernames":["jrpespinas"]}'
```

2. You should expect the following output on your terminal

```json
{
  "status": "success",
  "code": 200,
  "message": [
    {
      "name": "Bogs Espinas",
      "login": "jrpespinas",
      "company": "@demandscience",
      "followers": 8,
      "public_repos": 19
    }
  ]
}
```

## Trade offs

## References
