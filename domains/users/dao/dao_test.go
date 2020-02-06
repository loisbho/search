package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tj/assert"
	"search/domains/users"
	"testing"
)

var db *gorm.DB

func init() {
	v, err := gorm.Open("sqlite3", "tmp/gorm.db")
	if err != nil {
		panic(err)
	}
	db = v

	err = db.AutoMigrate(&User{}).Error
	if err != nil {
		panic(err)
	}
}

func TestDAO_ImportAll(t *testing.T) {
	ctx := context.Background()
	tx := db.Begin()
	defer tx.Rollback()
	dao := New(tx)

	//
	//Given
	u := []users.User{
		{
			ID:             1,
			Url:            "http://initech.zendesk.com/api/v2/users/1.json",
			ExternalID:     "74341f74-9c79-49d5-9611-87ef9b6eb75f",
			Name:           "Francisca Rasmussen",
			Alias:          "Miss Coffey",
			CreatedAt:      "2016-04-15T05:19:46 -10:00",
			Active:         true,
			Verified:       true,
			Shared:         false,
			Locale:         "en-AU",
			Timezone:       "Sri Lanka",
			LastLoginAt:    "2013-08-04T01:03:27 -10:00",
			Email:          "coffeyrasmussen@flotonic.com",
			Phone:          "8335-422-718",
			Signature:      "Don't Worry Be Happy!",
			OrganizationID: 119,
			Tags:           []string{"Springville", "Sutton", "Hartsville/Hartley", "Diaperville"},
			Suspended:      true,
			Role:           "admin",
		},
		{
			ID:             2,
			Url:            "http://initech.zendesk.com/api/v2/users/2.json",
			ExternalID:     "c9995ea4-ff72-46e0-ab77-dfe0ae1ef6c2",
			Name:           "Cross Barlow",
			Alias:          "Miss Joni",
			CreatedAt:      "2016-06-23T10:31:39 -10:00",
			Active:         true,
			Verified:       true,
			Shared:         false,
			Locale:         "en-AU",
			Timezone:       "Sri Lanka",
			LastLoginAt:    "2012-04-12T04:03:28 -10:00",
			Email:          "jonibarlow@flotonic.com",
			Phone:          "9575-552-585",
			Signature:      "Don't Worry Be Happy!",
			OrganizationID: 119,
			Tags:           []string{"Foxworth", "Woodlands", "Herlong", "Henrietta"},
			Suspended:      false,
			Role:           "admin",
		},
	}

	//When
	err := dao.ImportAll(ctx, users.ImportAllInput{Users: u})

	//Then
	assert.Nil(t, err)

	//verify record was stored
	got, err := dao.Find(ctx, users.FindAllInput{Field: "ID", Value: "1"})
	assert.Nil(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, "Francisca Rasmussen", got[0].Name)

	//verify record was stored
	got, err = dao.Find(ctx, users.FindAllInput{Field: "alias", Value: "Miss Joni"})
	assert.Nil(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, int32(2), got[0].ID)
}
