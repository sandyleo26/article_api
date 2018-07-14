Article API
--
## Setup
Before doing the following, You need to have docker and golang installed. `GOPATH` should be set correctly.


```
# clone repo
go get -u github.com/sandyleo26/article_api

# change to workdir
cd $GOPATH/src/github.com/sandyleo26/article_api

# start db container
make db-container-start

# db migrate
make db-up

# prepare build
make setup

# run
make
    
```


## Example
I use `curl` for testing. Other tool like postman will also work.

#### Using `curl`

```
# post an article
curl -X POST http://localhost:4321/articles -d "@article_data.json"

# get article by id
curl -i -X GET http://localhost:4321/articles/1

# get tag
curl -X GET http://localhost:4321/tags/abc/20180712
```
