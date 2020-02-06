package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tj/assert"
	"search/domains/tickets"
	"testing"
)

var db *gorm.DB

func init() {
	v, err := gorm.Open("sqlite3", "tmp/gorm.db")
	if err != nil {
		panic(err)
	}
	db = v

	err = db.AutoMigrate(&Ticket{}).Error
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
	inputTicket := []tickets.Ticket{
		{
			ID:             "436bf9b0-1147-4c0a-8439-6f79833bff5b",
			Url:            "http://initech.zendesk.com/api/v2/tickets/436bf9b0-1147-4c0a-8439-6f79833bff5b.json",
			ExternalID:     "9210cdc9-4bee-485f-a078-35396cd74063",
			CreatedAt:      "2016-04-28T11:19:34 -10:00",
			Type:           "incident",
			Subject:        "A Catastrophe in Korea (North)",
			Description:    "Nostrud ad sit velit cupidatat laboris ipsum nisi amet laboris ex exercitation amet et proident. Ipsum fugiat aute dolore tempor nostrud velit ipsum.",
			Priority:       "low",
			Status:         "pending",
			SubmitterID:    38,
			AssigneeID:      24,
			OrganizationID: 116,
			Tags:           []string{"Ohio", "Pennsylvania", "American Samoa", "Northern Mariana Islands"},
			HasIncidents:   false,
			DueAt:          "2016-07-31T02:37:50 -10:00",
			Via:            "web",
		}, {
			ID:             "674a19a1-c330-45fb-8b61-b4d77ba87130",
			Url:            "http://initech.zendesk.com/api/v2/tickets/674a19a1-c330-45fb-8b61-b4d77ba87130.json",
			ExternalID:     "050ea8ce-251c-44c8-b71c-535dd9072a74",
			CreatedAt:      "2016-03-07T08:24:53 -11:00",
			Type:           "task",
			Subject:        "A Drama in St. Pierre and Miquelon",
			Description:    "Incididunt exercitation voluptate eu laborum proident Lorem minim pariatur. Lorem culpa amet Lorem Lorem commodo anim deserunt do consectetur sunt.",
			Priority:       "low",
			Status:         "open",
			SubmitterID:    49,
			AssigneeID:      14,
			OrganizationID: 109,
			Tags:           []string{"Connecticut", "Arkansas", "Missouri", "Alabama"},
			HasIncidents:   false,
			DueAt:          "2016-08-15T06:13:11 -10:00",
			Via:            "voice",
		},
	}

	//When
	err := dao.ImportAll(ctx, tickets.ImportAllInput{Tickets: inputTicket})

	//Then
	assert.Nil(t, err)

	//verify record was stored
	got, err := dao.Find(ctx, tickets.FindAllInput{Field: "OrganizationID", Value: "116"})
	assert.Nil(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, "436bf9b0-1147-4c0a-8439-6f79833bff5b", got[0].ID)

	//verify record was stored
	got, err = dao.Find(ctx, tickets.FindAllInput{Field: "Priority", Value: "low"})
	assert.Nil(t, err)
	assert.Len(t, got, 2)

	//verify record was stored
	got, err = dao.Find(ctx, tickets.FindAllInput{Field: "type", Value: "task"})
	assert.Nil(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, "674a19a1-c330-45fb-8b61-b4d77ba87130", got[0].ID)
}
