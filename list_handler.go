package main

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/opengeo2016/listings/db"
	"github.com/paulmach/go.geojson"
)

// makeCriteria parses the query params from http request and constructs the
// criteria object to retrieve records from DB.
func makeCriteria(r *http.Request) *db.Criteria {
	// Parse params.
	q := r.URL.Query()
	c := db.Criteria{}

	if minPrice, err := strconv.Atoi(q.Get("min_price")); err == nil {
		c.MinPrice = &minPrice
	}

	if maxPrice, err := strconv.Atoi(q.Get("max_price")); err == nil {
		c.MaxPrice = &maxPrice
	}

	if minBed, err := strconv.Atoi(q.Get("min_bed")); err == nil {
		c.MinBed = &minBed
	}

	if maxBed, err := strconv.Atoi(q.Get("max_bed")); err == nil {
		c.MaxBed = &maxBed
	}

	if minBath, err := strconv.Atoi(q.Get("min_bath")); err == nil {
		c.MinBath = &minBath
	}

	if maxBath, err := strconv.Atoi(q.Get("max_bath")); err == nil {
		c.MaxBath = &maxBath
	}

	if startId, err := strconv.Atoi(q.Get("start_id")); err == nil {
		c.StartId = &startId
	}

	if numRequested, err := strconv.Atoi(q.Get("num_docs")); err == nil {
		// We are retrieving one more document than what user requested to determine if there is a "next page".
		numRequested = numRequested + 1
		c.NumRequested = &numRequested
	}

	return &c

}

// createNextUrl constructs the url for the next page.
func createNextUrl(url url.URL, startId int) string {
	q := url.Query()
	q.Set("start_id", strconv.Itoa(startId))
	url.RawQuery = q.Encode()

	return url.String()
}

// createFeature converts a ListingDetailsDocument into a geojson feature.
func createFeature(d db.ListingDetailsDocument) *geojson.Feature {
	f := geojson.NewPointFeature([]float64{d.Lng, d.Lat})
	f.SetProperty("id", d.RecId)
	f.SetProperty("street", d.Street)
	f.SetProperty("price", d.Price)
	f.SetProperty("bedrooms", d.NumOfBed)
	f.SetProperty("bathrooms", d.NumOfBath)
	f.SetProperty("sq_ft", d.SquareFeet)

	return f
}

// MakeListHandler creates a handler for proccesing listings request.
func MakeListHandler(lda *db.ListingDetailsAccessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := makeCriteria(r)
		documents, err := lda.List(c)

		if err != nil {
			log.Println("Failed in db retrieval. Url: ", r.URL.String(), " error: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fc := geojson.NewFeatureCollection()
		for idx, d := range *documents {
			if c.NumRequested == nil || idx < *c.NumRequested-1 {
				fc.AddFeature(createFeature(d))
			}
		}

		rawJSON, err := fc.MarshalJSON()
		if err != nil {
			log.Println("Failed in converting into json. Url: ", r.URL.String(), " error: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Sets next page link if there are more documents in the results than requested.
		if c.NumRequested != nil && len(*documents) >= *c.NumRequested {
			w.Header().Set("Link", createNextUrl(*r.URL, (*documents)[*c.NumRequested-1].RecId))
		}

		// Success.
		w.Write([]byte(rawJSON))
	}
}
