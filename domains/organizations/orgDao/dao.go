package orgDao

import (
	"context"
	"github.com/jinzhu/gorm"
	"golang.org/x/xerrors"
	"search/domains/organizations"
	"search/util"
	"search/validate"
	"time"
)

type DAO struct {
	db *gorm.DB
}

//New returns a new instance of DAO
func New(db *gorm.DB) *DAO {
	return &DAO{
		db: db,
	}
}

//Organization represents organization details
type Organization struct {
	ID            int32
	Url           string
	ExternalID    string
	Name          string
	DomainNames   string
	CreatedAt     time.Time
	Details       string
	SharedTickets bool
	Tags          string
}

//ImportAll stores the organizations into the database
func (d *DAO) ImportAll(ctx context.Context, input organizations.ImportAllInput) error {
	for _, item := range input.Orgs {
		org := Organization{
			ID:            item.ID,
			Url:           item.Url,
			ExternalID:    item.ExternalID,
			Name:          item.Name,
			DomainNames:   util.Join(item.DomainNames),
			CreatedAt:     util.ToTime(item.CreatedAt),
			Details:       item.Details,
			SharedTickets: item.SharedTickets,
			Tags:          util.Join(item.Tags),
		}

		tx := d.db.Create(&org)
		if err := tx.Error; err != nil {
			return err
		}
	}
	return nil
}

//Find searches by field and value
func (d *DAO) Find(ctx context.Context, input organizations.FindAllInput) ([]organizations.Organization, error) {
	if err := validate.Struct(input); err != nil {
		return nil, err
	}
	val := "%" + input.Value + "%"
	tx := d.db.Model(&Organization{}).
		Where(gorm.ToColumnName(input.Field)+" LIKE ?", val)

	var orgs []Organization
	if err := tx.Find(&orgs).Error; err != nil {
		return nil, xerrors.Errorf("unable to get organization for term: %s and value: %s: %w", input.Field, input.Value, err)
	}

	return newOrganizations(orgs)
}

func newOrganizations(orgs []Organization) ([]organizations.Organization, error) {
	result := make([]organizations.Organization, 0)
	for _, item := range orgs {
		result = append(result, organizations.Organization{
			ID:            item.ID,
			Url:           item.Url,
			ExternalID:    item.ExternalID,
			Name:          item.Name,
			DomainNames:   util.Split(item.DomainNames),
			CreatedAt:     item.CreatedAt.String(),
			Details:       item.Details,
			SharedTickets: item.SharedTickets,
			Tags:          util.Split(item.Tags),
		})
	}
	return result, nil
}
