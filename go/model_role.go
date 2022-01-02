package openapi

// Role - A role struct to cater for different roles, e.g. gigworker, employer etc
type Role struct {
	Id string `bson:"id,omitempty"`

	Name string `bson:"name,omitempty"`
}
