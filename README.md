Article API
--
## Setup
Before doing the following, You need to have docker and golang installed. `GOPATH` should be set correctly.


```
# clone repo
go get -u github.com/sandyleo26/article_api

# prepare build
make setup

# change to workdir
cd $GOPATH/src/github.com/sandyleo26/article_api

# start db container
make db-container-start

# db migrate
make db-up

# run
make
    
```

## Example
I use `curl` for testing. Other tool like postman will also work.

```
# post an article
curl -X POST http://localhost:4321/articles -d "@example/article_data.json"

# get article by id
curl -i -X GET http://localhost:4321/articles/1

# get tag
curl -X GET http://localhost:4321/tags/abc/20180712
```

## Solution description
In simple words, the api service take HTTP requests and write or retrieve data from database (postgres in a container) and then send back json response.

First, I start a web server to listen on port 4321. I use `mux` library to define routes and corresponding request handlers.

Then handler decode the request, parse and validate path parameters if exist and then send to different usecase functions.

The usecase functions are main business logic and iteract with database through ORM library `gorm`.

For `GetTag` usecase (`AddArticle` and `GetArticle` are straightforward)
  1. query articles within 24 hours of the date ordered by time of creating
  2. collect the articles that contains the specified tag
  3. collect related tags except the chosen one
  4. collect the 10 most recent article IDs

Finally, the tests are on all on usecase level, asserting results according to the specification.

## Assumptions
* Databases are gorm compatable
* It has no authurization or rate limiting which means everyone is allowed to read and write database without limit
* Query time is local time AEST
* Limited of string length imposed by database

## Thought
I took me 2h to finish the initial prototype. Then another 1h to complete and refine. And 0.5h for writing tests and 1h to finish the README and final test. So in total 4.5h.

I think it's a typical REST service that represent many common problems such routing, data persistence, authorization etc. Things I'd try if time is permitted:

* Add authorization
* Add a front end or return templated html
* Try GraphQL
* Add testing for handlers
