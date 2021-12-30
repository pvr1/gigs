package openapi

// Gig - A gig struct
type Gig struct {
	Id string `json:"id,omitempty"`

	Category Category `json:"category,omitempty"`

	Name string `json:"name"`

	Description []string `json:"description"`

	Measurableoutcome []string `json:"measurableoutcome"`

	Tags []Tag `json:"tags,omitempty"`

	// gig status in the store
	Status string `json:"status,omitempty"`
}

var gigs = []Gig{
	{Id: "1", Name: "Gig 1", Description: []string{"description 1"}, Measurableoutcome: []string{"measurableoutcome 1"}, Status: "available"},
	{Id: "2", Name: "Gig 2", Description: []string{"description 2"}, Measurableoutcome: []string{"measurableoutcome 2"}, Status: "available"},
	{Id: "3", Name: "Gig 3", Description: []string{"description 3"}, Measurableoutcome: []string{"measurableoutcome 3"}, Status: "available"},
	{Id: "4", Name: "Gig 3", Description: []string{"description 3"}, Measurableoutcome: []string{"measurableoutcome 3"}, Status: "sold"},
	{Id: "5", Name: "Gig 3", Description: []string{"description 3"}, Measurableoutcome: []string{"measurableoutcome 3"}, Status: "pending"},
}
