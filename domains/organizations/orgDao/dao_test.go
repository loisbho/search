package orgDao

import (
	"context"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tj/assert"
	"search/domains/organizations"
	"testing"
)

var db *gorm.DB

func init() {
	v, err := gorm.Open("sqlite3", "tmp/gorm.db")
	if err != nil {
		panic(err)
	}
	db = v

	err = db.AutoMigrate(&Organization{}).Error
	if err != nil {
		panic(err)
	}
}

func TestDAO_ImportAll(t *testing.T) {
	ctx := context.Background()
	tx := db.Begin()
	defer tx.Rollback()

	dao := New(tx)
	//Given
	orgs := []organizations.Organization{
		{
			ID:   1,
			Name: "Sulfax",
			DomainNames: []string{
				"comvey.com",
				"velity.com",
				"enormo.com",
			},
			CreatedAt:     "2016-01-12T01:16:06 -11:00",
			Details:       "MegaCörp",
			SharedTickets: true,
			Tags: []string{
				"Travis",
				"Clarke",
				"Glenn",
				"Santos",
			},
		},
		{
			ID:   2,
			Name: "org2",
		},
	}

	//When
	err := dao.ImportAll(ctx, organizations.ImportAllInput{Orgs: orgs})

	//Then
	assert.Nil(t, err)

	//verify record was stored
	got, err := dao.Find(ctx, organizations.FindAllInput{Field: "Name", Value: "Sulfax"})
	assert.Nil(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, []string{"comvey.com", "velity.com", "enormo.com"}, got[0].DomainNames)

	//verify record was stored
	got, err = dao.Find(ctx, organizations.FindAllInput{Field: "ID", Value: "2"})
	assert.Nil(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, "org2", got[0].Name)
}
