package openapi

// Gig - A gig struct
type Gig struct {
	Id string `bson:"id,omitempty"`

	Category Category `bson:"category,omitempty"`

	Name string `bson:"name"`

	Description []string `bson:"description"`

	Measurableoutcome []string `bson:"measurableoutcome"`

	Tags []Tag `bson:"tags,omitempty"`

	// gig status in the store
	Status string `bson:"status,omitempty"`
}

var gigs = []Gig{
	{Id: "1", Name: "Gig 1", Description: []string{"description 1"}, Measurableoutcome: []string{"measurableoutcome 1"}, Status: "available"},
	{Id: "2", Name: "Gig 2", Description: []string{"description 2"}, Measurableoutcome: []string{"measurableoutcome 2"}, Status: "available"},
	{Id: "3", Name: "Gig 3", Description: []string{"description 3"}, Measurableoutcome: []string{"measurableoutcome 3"}, Status: "available"},
	{Id: "4", Name: "Gig 3", Description: []string{"description 3"}, Measurableoutcome: []string{"measurableoutcome 3"}, Status: "sold"},
	{Id: "5", Name: "Gig 3", Description: []string{"description 3"}, Measurableoutcome: []string{"measurableoutcome 3"}, Status: "pending"},
}
