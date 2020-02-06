package users

import "context"

//User is a representation of a user
type User struct {
	ID             int32    `json:"_id"`
	Url            string   `json:"url"`
	ExternalID     string   `json:"external_id"`
	Name           string   `json:"name"`
	Alias          string   `json:"alias"`
	CreatedAt      string   `json:"created_at"`
	Active         bool     `json:"active"`
	Verified       bool     `json:"verified"`
	Shared         bool     `json:"shared"`
	Locale         string   `json:"locale"`
	Timezone       string   `json:"timezone"`
	LastLoginAt    string   `json:"last_login_at"`
	Email          string   `json:"email"`
	Phone          string   `json:"phone"`
	Signature      string   `json:"signature"`
	OrganizationID int32    `json:"organization_id"`
	Tags           []string `json:"tags"`
	Suspended      bool     `json:"suspended"`
	Role           string   `json:"role"`
}

//FindAllInput fields necessary for Find
type FindAllInput struct {
	Field string `validate:"required"`
	Value string `validate:"required"`
}

//ImportAllInput includes the users that will be imported
type ImportAllInput struct {
	Users []User
}

type Finder interface {
	//Find returns the record for the field
	Find(ctx context.Context, input FindAllInput) ([]User, error)
}

type Mutator interface {
	//ImportAll updates the database
	ImportAll(ctx context.Context, input ImportAllInput) error
}
