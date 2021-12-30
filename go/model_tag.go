package openapi

type Tag struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

var tags = []Tag{
	{Id: 1, Name: "tag"},
	{Id: 2, Name: "1"},
	{Id: 3, Name: "2"},
}
