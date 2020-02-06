package organizations

import "context"

//Organization is a representation of an organization
type Organization struct {
	ID            int32    `json:"_id"`
	Url           string   `json:"url"`
	ExternalID    string   `json:"external_id"`
	Name          string   `json:"name"`
	DomainNames   []string `json:"domain_names"`
	CreatedAt     string   `json:"created_at"`
	Details       string   `json:"details"`
	SharedTickets bool     `json:"shared_tickets"`
	Tags          []string `json:"tags"`
}

//FindAllInput required fields
type FindAllInput struct {
	Field string `validate:"required"`
	Value string `validate:"required"`
}

//ImportAllInput
type ImportAllInput struct {
	Orgs []Organization
}

//Finder queries the database
type Finder interface {
	//Find returns the record for the field
	Find(ctx context.Context, input FindAllInput) ([]Organization, error)
}

//Mutator modifies the database
type Mutator interface {
	//ImportAll updates the database
	ImportAll(ctx context.Context, input ImportAllInput) error
}
