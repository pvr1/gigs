package openapi

// User - A user struct used to store the user information
type User struct {
	Id int64 `bson:"id,omitempty"`

	Username string `bson:"username,omitempty"`

	FirstName string `bson:"firstName,omitempty"`

	LastName string `bson:"lastName,omitempty"`

	Email string `bson:"email,omitempty"`

	Password string `bson:"password,omitempty"`

	Phone string `bson:"phone,omitempty"`

	// User Status
	UserStatus int32 `bson:"userStatus,omitempty"`

	// User Role - e.g. gigworker, employer etc
	Role []Role `bson:"role,omitempty"`
}
