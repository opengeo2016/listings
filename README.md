# listings
## Overview
This repo implements a simple server to list and filter a number of houses based on certain criteria, such as pricing. The results are presented in geojson format and could be visualized directly in [geojson.io](http://geojson.io/).

## Technical stack
- go server
- mongo db
- [mgo](https://github.com/go-mgo/mgo) for db connection
- [go.geojson](https://github.com/paulmach/go.geojson) for recordign result.

## File description
- server.go includes the logic to setup the server.
- list_handler.go includes the logic to process listing requests from user.
- db/listing_details.go specifies the schema of data stored in mongodb.
- db/listing_details_accessor.go includes the logic to access and filter data from mongodb.

## How to run from scratch
### 1. Set up [mongodb](https://docs.mongodb.org/manual/installation/)
### 2. Populate data
- populate data into mongo db using the following command.
```
mongoimport -d opendoor -c ListingDetails --type csv --file listing-details.csv --headerline
```
- Note:
- We are specifying the mongo db name to be "opendoor", and collection name to be "ListingDetails".
- listing-details.csv is the data file containing all the housing data. Its heading must include the following fields
|id|street|status|price|bedrooms|bathrooms|sq_ft|lat|lng|, and here is one example record

| id | street | status | price | bedrooms | bathrooms | sq_ft | lat | lng|
|----|--------|--------|-------|----------|-----------|-------|-----|----|
| 0  |545 2nd Pl|	pending|299727|4|1|1608|33.3694442|-112.1197146|

### 3. Build server
- Install [golang](https://golang.org/doc/install)
- checkout code to $GOPATH/src/github.com/opengeo2016/listings
```
  mkdir -p $GOPATH/src/github.com/opengeo2016
  cd $GOPATH/src/github.com/opengeo2016
  git clone https://github.com/opengeo2016/listings
```
- Install deps
```
  go get github.com/paulmach/go.geojson
  go get gopkg.in/mgo.v2
```
- Build and run
```
  cd $GOPATH/src/github.com/opengeo2016/listings
  go build .
  ./listings --db_addr=${connection to mongodb} --service_port=${port to serve request}
```

## Api spec
The server supports a single api with different params.
```
/listings
```
### Input query params:
- min_price, the minimal price to filter based on. 
- max_price, the maximal price to filter based on. 
- min_bed, the minimal number of bedrooms to filter based on. 
- max_bed, the maximal number of bedrooms to filter based on. 
- min_bath, the minimal number of bathrooms to filter based on.
- max_bath, the maximal number of bathrooms to filter based on.
- num_docs, the maximal number of record to be returned.
- start_id, the minimal id of the record to start the search with. The id of the returned records should be greater or equal to it.
E.g. The following request retrieves 10 listings with price being equal or higher than 115254, starting from the record with id being 1001.
```
/listings?min_price=115254&num_docs=10&start_id=1001
```

### Output
- The returned results are sorted by their ids from lowest to highest, and are encoded in geojson format.
- If there are more results than requested by the "num_doc" param, an additional "Link" field will be provided in the response header, which specifies the url to the next set of data. E.g. the following could be included in the response header of the above request.
```
Link: /listings?min_price=115254&num_docs=10&start_id=1012
```

## TODO
- Improve code quality: getting rid of hard-coded strings, add unit tests.
- API improvement: handle CORS, returning more detailed error message when there is an error.
- db improvement: creating indices to improve performance.
- supporting spatial queries: user are allowed to specify a region (rectangle or circle drawed on map, or input a city/county name), return filter results within the specified area.
- supporting summary: in the returning data set, provides summary, such as min/max/avg, of the resulting data.
- supporting hisotry: store and enable searches on price history.
- supporting more pagniation operations: to the first page, to the last page, to a specified page, ...
- and more ...

