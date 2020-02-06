package tickets

import "context"

//Ticket is a representation of a ticket
type Ticket struct {
	ID             string   `json:"_id"`
	Url            string   `json:"url"`
	ExternalID     string   `json:"external_id"`
	CreatedAt      string   `json:"created_at"`
	Type           string   `json:"type"`
	Subject        string   `json:"subject"`
	Description    string   `json:"description"`
	Priority       string   `json:"priority"`
	Status         string   `json:"status"`
	SubmitterID    int32    `json:"submitter_id"`
	AssigneeID      int32    `json:"assignee_id"`
	OrganizationID int32    `json:"organization_id"`
	Tags           []string `json:"tags"`
	HasIncidents   bool     `json:"has_incidents"`
	DueAt          string   `json:"due_at"`
	Via            string   `json:"via"`
}

//FindAllInput fields necessary to Find
type FindAllInput struct {
	Field string `validate:"required"`
	Value string `validate:"required"`
}

//ImportAllInput includes the users that will be imported
type ImportAllInput struct {
	Tickets []Ticket
}

type Finder interface {
	//Find returns the record(s) found
	Find(ctx context.Context, input FindAllInput) ([]Ticket, error)
}

type Mutator interface {
	//ImportAll updates the database
	ImportAll(ctx context.Context, input ImportAllInput) error
}
