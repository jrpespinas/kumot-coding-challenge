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

## Solution

### Solution to the Objectives

1. I made a switch statement that takes in a string to determine which field of the `domain.User` to sort. Each case assigns a function to the variable `less` which will then be used by the `sort.Slice` function. [[5]](#5)

2. I looped through the list of usernames and append found users to `[]domain.User`. If the username does not exist, then we can just simply skip the current iteration.

3. Use redis and implement two major functions `Set` and `Get`. We set the expiration of a key-value pair to 2 minutes in the `Set` function. If the user's information has been cached, we simply use the `Get` function to retrieve the data.

4. User friendly responses such as "service is unavailable" were used instead of programming jargons "ERROR: unable to marshal json".

5. Was used `https://api.github.com/users/{username}`, not 3rd party libraries.

### Deployment: Docker

Using `docker-compose` deploying both the redis server and the golang server is easy. By cloning this repository, you will be able to observe the same performance and correctness of the application with any compatibility issues as long as you have docker installed.

The following is a rough summary of the docker compose which represents the architecture of the system.

```Shell
.
└─ services
   ├── server
   └── redis-server
```

### Architecture: Clean Architecture

Clean Architecture[[1]](#1) is a system architecture guideline proposed by Robert C. Martin inspired by similar architecture like hexagonal and onion architecture. To put it simply, the main point of this approach is the separation of concerns. The dependency points deeper in the application. Meaning the outermost layer must only depend from the inner layer, rather than depending on external entities. This approach enables the application to be testable, maintainable, changeable, independent of frameworks, and easy to deploy.

To elaborate, there are **three layers** in this architecture: repository, service, and controller.

The **Controller** layer is responsible for handling the request and for returning the responses. It needs the service layer to function effectively, without the service layer the controller layer will not be able to process requests and return responses.

The **Service** layer is responsible for the business logic. This is the layer in which you manipulate the domain/models. In this case, this is where we validate the request, generate token, sort the retrieved data from github, pulling and setting data into the cache. The service layer needs the repository layer to be able to manipulate domain/models.

The **Repository** layer is responsible for the data access or the data persistence layer. This is where we store the key-value pair in our cache. This is where we pull data from. In this case, where access to github api is implemented.

### Project Structure: Domain Driven Design

Clean Architecture works really well with Domain-Driven Design (DDD). [[2]](#2)[[3]](#3)

To take a well explained quote from StackOverflow:

> DDD is about trying to make your software a model of a real-world system or process. In using DDD, you are meant to work closely with a domain expert who can explain how the real-world system works. For example, if you're developing a system that handles the placing of bets on horse races, your domain expert might be an experienced bookmaker.
> Between yourself and the domain expert, you build a ubiquitous language (UL), which is basically a conceptual description of the system. The idea is that you should be able to write down what the system does in a way that the domain expert can read it and verify that it is correct. In our betting example, the ubiquitous language would include the definition of words such as 'race', 'bet', 'odds' and so on.
> [[4]](#4)

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
│   │       ├── token.go             # Struct definition for token response
│   │       └── user.go              # Struct definition for parsing request body
│   │
│   ├── listing                      # Service Layer: display of github details
│   │   ├── service.go               # Handle caching and repository calls
│   │   ├── service_test.go          # Unit and mock test for the service layer
│   │   └── sorter.go                # Define helper function for sorting the github details
│   │
│   ├── logging                      # Service Layer: logging
│   │   └── logger.go                # Uses zerolog as the logger
│   │
│   ├── repository                   # Repository Layer: Data Access and Persistence
│   │   ├── api
│   │   │   └── repository.go        # Definition of Repository interface, github implementation
│   │   └── cache
│   │       └── cache.go             # Definition of Cache Interface, Redis Implementation
│   │
│   ├── router
│   │   └── router.go                # Definition of router interface, Chi-Router implementation
│   │
│   └── session                      # Service Layer: Session
│       └── service.go               # Generate and Verify token using redis
└── redis.conf                       # Basic custom configuration for our redis Cache
```

## Installation and Usage

### Setup

Make sure you have Docker installed on your device, otherwise you may [download](https://www.docker.com/products/docker-desktop/) docker from their webpage. Installing docker desktop also includes Docker CLI with Docker Compose plugin.

If you already have Docker, simply follow the steps to deploy our backend and redis.

1. Open your linux terminal and clone this repository.

```shell
$ git clone git@github.com:jrpespinas/kumot-coding-challenge.git
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

1. Open another linux terminal and paste the following curl request to generate your token

```shell
$ curl -X GET http://localhost:8080/generate-token -H 'Content-Type: application/json'
```

2. You will receive your token in this format

```json
{
  "status": "success",
  "code": 200,
  "data": { "token": "8653dd19ea99820b91a5e5815bf5b10b" }
}
```

3. Next, copy the token (without the quotation marks) and insert it in your header on the following request

```shell
$ curl -X POST http://localhost:8080/users -H 'Content-Type: application/json' -H 'Session-Token: <INSERT-TOKEN-HERE>' -d '{"usernames":["jrpespinas"]}'
```

See this example

```shell
$ curl -X POST http://localhost:8080/users -H 'Content-Type: application/json' -H 'Session-Token: 8653dd19ea99820b91a5e5815bf5b10b' -d '{"usernames":["jrpespinas"]}'
```

4. You should expect the following output on your terminal

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

1. Ideally, Clean Architecture is used for large enterprise applications and microservices. The reason I chose this architecture is inspired by the fact that your company uses microservices.

2. Another trade off for Clean Architecture is the increased initial development time. Writing more code at the start to compensate for better maintainability, testability, and minimal refactoring.

3. Domain-Driven Design is the appropriate approach if I was working domain expert who can explain how the real-world system works and who must be able to understand the project structure. The reason for this approach is for explanation purposes I consider you as domain expert who needs to understand my project structure.

4. I used a primitive approach by using rand in generating tokens to utilize Redis. Better to use other methods such as JWT.

## References

<a id="1">[1]</a> [Clean Architecture](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164)

<a id="2">[2]</a> [Domain-Driven Design: Tackling Complexity in the Heart of Software](https://www.amazon.com/gp/product/0321125215?ie=UTF8&linkCode=sl1&tag=&linkId=f5b14b978ccb68e3d1c34eefcfe9ff99&language=en_US&ref_=as_li_ss_tl)

<a id="3">[3]</a> [How Do You Structure Your Go Apps](https://www.youtube.com/watch?v=oL6JBUk6tj0)

<a id="4">[4]</a> [What is Domain Driven Design (DDD)?](https://stackoverflow.com/a/1222488)

<a id="5">[5]</a> [Sort Struct Slice by Specific Field](https://stackoverflow.com/a/53025566)
