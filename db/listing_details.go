package db

// ListingDetailsDocument specifies the document type in db collection.
type ListingDetailsDocument struct {
	// The id defined by data provider, which is different from '_id' specified by mongodb.
	RecId int `bson:"id,omitempty"`

	// The street address.
	Street string `bson:"street,omitempty"`

	// TODO: put some restictions on possible values, instead of an arbitrary string.
	// The status of the listing. (pending/active/sold etc.)
	Status string `bson:"status,omitempty"`

	// The price, # of bedrooms, # of bath rooms, square feet of the listing.
	Price      int `bson:"price,omitempty"`
	NumOfBed   int `bson:"bedrooms,omitempty"`
	NumOfBath  int `bson:"bathrooms,omitempty"`
	SquareFeet int `bson:"sq_ft,omitempty"`

	// The latitude/longitude of the listing.
	Lat float64 `bson:"lat,omitempty"`
	Lng float64 `bson:"lng,omitempty"`
}
