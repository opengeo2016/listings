package db

import (
	"math"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB_NAME    = "opendoor"
	COLLECTION = "ListingDetails"
)

// The criteria to specify when retrieving listings from db.
// Setting nil if a criteria should be ignored.
type Criteria struct {
	// Price criteria.
	MinPrice *int
	MaxPrice *int

	// Num of bedrooms criteria.
	MinBed *int
	MaxBed *int

	// Num of bathrooms criteria.
	MinBath *int
	MaxBath *int

	// Document Id to start with.
	// Id of all the documents in the result should be greater or equal to StartId.
	StartId *int

	// Number of documents to retrieve.
	NumRequested *int
}

type ListingDetailsAccessor struct {
	// Collection.
	c *mgo.Collection
}

// NewListingDetailsAccessor creates a ListingDetailsAccessor.
func NewListingDetailsAccessor(s *mgo.Session) *ListingDetailsAccessor {
	return &ListingDetailsAccessor{
		c: s.DB(DB_NAME).C(COLLECTION),
	}
}

// constructs a filter on a field to access mongodb.
func getFilterSpec(minV *int, maxV *int) *bson.M {
	if minV == nil && maxV == nil {
		return nil
	}
	spec := bson.M{}

	if minV != nil {
		spec["$gte"] = *minV
	}

	if maxV != nil {
		spec["$lte"] = *maxV
	}
	return &spec
}

// List accessed the db collection and returns all records based on Criteria specified in the input object.
func (lda *ListingDetailsAccessor) List(c *Criteria) (documents *[]ListingDetailsDocument, e error) {
	filter := bson.M{}

	if spec := getFilterSpec(c.MinPrice, c.MaxPrice); spec != nil {
		filter["price"] = spec
	}

	if spec := getFilterSpec(c.MinBed, c.MaxBed); spec != nil {
		filter["bedrooms"] = spec
	}

	if spec := getFilterSpec(c.MinBath, c.MaxBath); spec != nil {
		filter["bathrooms"] = spec
	}

	if spec := getFilterSpec(c.StartId, nil); spec != nil {
		filter["id"] = spec
	}

	limit := math.MaxInt32
	if c.NumRequested != nil {
		limit = *c.NumRequested
	}

	iter := lda.c.Find(filter).Sort("id").Limit(limit).Iter()

	var result []ListingDetailsDocument
	err := iter.All(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
