package openapi

// Tag - A tag struct used to meta tag the gigs
type Tag struct {
	Id   int64  `bson:"id,omitempty"`
	Name string `bson:"name,omitempty"`
}

// The possible values for the meta tag
var tags = []Tag{
	{Id: 1, Name: "tag"},
	{Id: 2, Name: "1"},
	{Id: 3, Name: "2"},
}
