# Kumot Coding Challenge

## Problem Statement

Golang project that has an API endpoint that takes a list of github
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

## Installation

## Architecture and Design

## Trade offs

## Other references
