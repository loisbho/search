package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"golang.org/x/xerrors"
	"search/domains/users"
	"search/util"
	"time"
)

type DAO struct {
	db *gorm.DB
}

func New(db *gorm.DB) *DAO {
	return &DAO{
		db: db,
	}
}

//User is a representation of a user
type User struct {
	ID             int32
	Url            string
	ExternalID     string
	Name           string
	Alias          string
	CreatedAt      time.Time
	Active         bool
	Verified       bool
	Shared         bool
	Locale         string
	Timezone       string
	LastLoginAt    time.Time
	Email          string
	Phone          string
	Signature      string
	OrganizationID int32
	Tags           string
	Suspended      bool
	Role           string
}

func (d *DAO) ImportAll(ctx context.Context, input users.ImportAllInput) error {
	for _, item := range input.Users {
		u := User{
			ID:             item.ID,
			Url:            item.Url,
			ExternalID:     item.ExternalID,
			Name:           item.Name,
			Alias:          item.Alias,
			CreatedAt:      util.ToTime(item.CreatedAt),
			Active:         item.Active,
			Verified:       item.Verified,
			Shared:         item.Shared,
			Locale:         item.Locale,
			Timezone:       item.Timezone,
			LastLoginAt:    util.ToTime(item.LastLoginAt),
			Email:          item.Email,
			Phone:          item.Phone,
			Signature:      item.Signature,
			OrganizationID: item.OrganizationID,
			Tags:           util.Join(item.Tags),
			Suspended:      item.Suspended,
			Role:           item.Role,
		}
		tx := d.db.Create(&u)
		if err := tx.Error; err != nil {
			//todo idempotent, if record already exists, continue
			return err
		}
	}
	return nil
}

func (d *DAO) Find(ctx context.Context, input users.FindAllInput) ([]users.User, error) {
	//todo add validation
	tx := d.db.Model(&User{}).
		Where(gorm.ToColumnName(input.Field)+"= ?", input.Value)

	var users []User
	if err := tx.Find(&users).Error; err != nil {
		return nil, xerrors.Errorf("unable to get user for term: %s and value: %s: %w", input.Field, input.Value, err)
	}

	return newUsers(users)
}

func newUsers(u []User) ([]users.User, error) {
	result := make([]users.User, 0)
	for _, item := range u {
		result = append(result, users.User{
			ID:             item.ID,
			Url:            item.Url,
			ExternalID:     item.ExternalID,
			Name:           item.Name,
			Alias:          item.Alias,
			CreatedAt:      item.CreatedAt.String(),
			Active:         item.Active,
			Verified:       item.Verified,
			Shared:         item.Shared,
			Locale:         item.Locale,
			Timezone:       item.Timezone,
			LastLoginAt:    item.LastLoginAt.String(),
			Email:          item.Email,
			Phone:          item.Phone,
			Signature:      item.Signature,
			OrganizationID: item.OrganizationID,
			Tags:           util.Split(item.Tags),
			Suspended:      item.Suspended,
			Role:           item.Role,
		})
	}
	return result, nil
}
