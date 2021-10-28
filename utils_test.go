package mongo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func TestFlattenedMapFromInterface(t *testing.T) {
	type bar struct {
		Action string `json:"action"`
	}
	type dummy struct {
		Name *string  `json:"name"`
		Bar  bar      `json:"bar"`
		Bars []string `json:"bars"`
	}

	dummyName := "dummy"

	tests := []struct {
		dummy       dummy
		expectedRes map[string]interface{}
	}{
		{
			dummy: dummy{
				Name: &dummyName,
				Bar: bar{
					Action: "dumb",
				},
				Bars: nil,
			},
			expectedRes: map[string]interface{}{
				"name":       "dummy",
				"bar.action": "dumb",
			},
		},
		{
			dummy: dummy{
				Bar: bar{
					Action: "dumb",
				},
				Bars: []string{"tata", "toto"},
			},
			expectedRes: map[string]interface{}{
				"bar.action": "dumb",
				"bars":       []interface{}{"tata", "toto"},
			},
		},
		{
			dummy: dummy{
				Name: &dummyName,
				Bar: bar{
					Action: "dumb",
				},
				Bars: []string{"tata", "toto"},
			},
			expectedRes: map[string]interface{}{
				"name":       "dummy",
				"bar.action": "dumb",
				"bars":       []interface{}{"tata", "toto"},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedRes, FlattenedMapFromInterface(test.dummy))
	}
}

func TestFlattenedMapFromBson(t *testing.T) {
	obj := bson.M{"test": "foo", "test2": nil}

	assert.Equal(t, map[string]interface{}{"test": "foo", "test2": nil}, FlattenedMapFromInterface(obj))
}

type dummyDateObj struct {
	BasicDocument    `bson:",inline"`
	Start            time.Time
	End              time.Time
	DummyStupidField string
}

func (d dummyDateObj) DocumentName() string {
	return "dummyDate"
}

func TestFlattenedMapFromDates(t *testing.T) {
	client, _ := dummyConnect()

	obj := dummyDateObj{
		Start:            time.Now(),
		End:              time.Now().Add(24 * time.Hour),
		DummyStupidField: "pato",
	}

	client.Persist(&obj)

	obj2 := dummyDateObj{
		Start:            time.Now().Add(24 * time.Hour),
		End:              time.Now().Add(72 * time.Hour),
		DummyStupidField: "pata",
	}

	client.Update(&obj, obj.ID, obj2)

	assert.NotNil(t, obj)
}

type User struct {
	BasicDocument     `bson:",inline"`
	FirstName         string           `json:"firstName" bson:"firstName"`
	LastName          string           `json:"lastName" bson:"lastName"`
	Email             string           `json:"email"`
	ProfileCompletion int32            `json:"profileCompletion" bson:"profileCompletion"`
	Info              UserInfo         `json:"info,omitempty"`
	EmergencyContact  EmergencyContact `json:"emergencyContact,omitempty" bson:"emergencyContact"`
}

type UserInfo struct {
	Picture       string   `json:"picture"`
	Title         string   `json:"title"`
	Address1      string   `json:"address1" bson:"address1"`
	Address2      string   `json:"address2" bson:"address2"`
	City          string   `json:"city"`
	State         string   `json:"state"`
	Zip           string   `json:"zip"`
	Phone1        string   `json:"phone1" bson:"phone1"`
	Phone2        string   `json:"phone2" bson:"phone2"`
	ContactMethod []string `json:"contactMethod" bson:"contactMethod"`
}

type EmergencyContact struct {
	FirstName    string `json:"firstName" bson:"firstName"`
	LastName     string `json:"lastName" bson:"lastName"`
	Relationship string `json:"relationship"`
	Phone1       string `json:"phone1" bson:"phone1"`
	Phone2       string `json:"phone2" bson:"phone2"`
}

func TestFlattenedMapOnUser(t *testing.T) {
	userToUpdate := User{
		BasicDocument:     BasicDocument{},
		FirstName:         "Alex",
		LastName:          "Khoury",
		Email:             "alex@email.com",
		ProfileCompletion: 3,
		Info: UserInfo{
			Picture:       "qwc",
			Title:         "Mr",
			Address1:      "as",
			Address2:      "as",
			City:          "as",
			State:         "as",
			Zip:           "12345",
			Phone1:        "1111111111",
			Phone2:        "1111111111",
			ContactMethod: []string{"email", "sms"},
		},
		EmergencyContact: EmergencyContact{
			FirstName:    "as",
			LastName:     "as",
			Relationship: "as",
			Phone1:       "1111111111",
			Phone2:       "1111111111",
		},
	}

	res := FlattenedMapFromInterface(userToUpdate)

	for k, v := range res {
		fmt.Println(fmt.Sprintf("KEY %s: %s", k, v))
	}
}
