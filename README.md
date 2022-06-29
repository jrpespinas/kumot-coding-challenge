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

## Trade offs

## Other references
