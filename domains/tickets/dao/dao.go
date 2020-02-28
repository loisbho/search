package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"golang.org/x/xerrors"
	"search/domains/tickets"
	"search/util"
	"search/validate"
	"time"
)

type DAO struct {
	db *gorm.DB
}

//New returns a new DAO
func New(db *gorm.DB) *DAO {
	return &DAO{
		db: db,
	}
}

//Ticket is a representation of a ticket
type Ticket struct {
	ID             string
	Url            string
	ExternalID     string
	CreatedAt      time.Time
	Type           string
	Subject        string
	Description    string
	Priority       string
	Status         string
	SubmitterID    int32
	AsigneeID      int32
	OrganizationID int32
	Tags           string
	HasIncidents   bool
	DueAt          time.Time
	Via            string
}

//ImportAll stores tickets into database
func (d *DAO) ImportAll(ctx context.Context, input tickets.ImportAllInput) error {
	for _, item := range input.Tickets {
		t := Ticket{
			ID:             item.ID,
			Url:            item.Url,
			ExternalID:     item.ExternalID,
			CreatedAt:      util.ToTime(item.CreatedAt),
			Type:           item.Type,
			Subject:        item.Subject,
			Description:    item.Description,
			Priority:       item.Priority,
			Status:         item.Status,
			SubmitterID:    item.SubmitterID,
			AsigneeID:      item.AssigneeID,
			OrganizationID: item.OrganizationID,
			Tags:           util.Join(item.Tags),
			HasIncidents:   item.HasIncidents,
			DueAt:          util.ToTime(item.DueAt),
			Via:            item.Via,
		}

		tx := d.db.Create(&t)
		if err := tx.Error; err != nil {
			//todo idempotent, if record already exists, continue
			return err
		}
	}
	return nil
}

//Find searches by field and value
func (d *DAO) Find(ctx context.Context, input tickets.FindAllInput) ([]tickets.Ticket, error) {
	if err := validate.Struct(input); err != nil {
		return nil, err
	}
	val := "%" + input.Value + "%"
	tx := d.db.Model(&Ticket{}).
		Where(gorm.ToColumnName(input.Field)+" LIKE ?", val)

	var tickets []Ticket
	if err := tx.Find(&tickets).Error; err != nil {
		return nil, xerrors.Errorf("unable to get tickets for term: %s and value: %s: %w", input.Field, input.Value, err)
	}

	return newTickets(tickets)
}

func newTickets(t []Ticket) ([]tickets.Ticket, error) {
	result := make([]tickets.Ticket, 0)
	for _, item := range t {
		result = append(result, tickets.Ticket{
			ID:             item.ID,
			Url:            item.Url,
			ExternalID:     item.ExternalID,
			CreatedAt:      item.CreatedAt.String(),
			Type:           item.Type,
			Subject:        item.Subject,
			Description:    item.Description,
			Priority:       item.Priority,
			Status:         item.Status,
			SubmitterID:    item.SubmitterID,
			AssigneeID:      item.AsigneeID,
			OrganizationID: item.OrganizationID,
			Tags:           util.Split(item.Tags),
			HasIncidents:   item.HasIncidents,
			DueAt:          item.DueAt.String(),
			Via:            item.Via,
		})
	}
	return result, nil
}
