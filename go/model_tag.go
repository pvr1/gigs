package openapi

// Tag - A tag struct used to meta tag the gigs
type Tag struct {
	Id   int64  `bson:"id,omitempty"`
	Name string `bson:"name,omitempty"`
}
