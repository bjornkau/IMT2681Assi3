package main

import (
	"gopkg.in/mgo.v2"
	"testing"
	//"gopkg.in/mgo.v2/bson"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
)

type Result struct {
	Parameters map[string]string `json:"parameters"`
}

type TestInput struct {
	Result2 Result `json:"result"`
}

//setting up the database for testing
func setUpDB(t *testing.T) *APIMongoDB {
	db := APIMongoDB{
		"mongodb://localhost:27017",
		"cloudtesting",
		"fixers",
	}

	session, err := mgo.Dial(db.Host)
	defer session.Close()

	if err != nil {
		t.Error(err)
	}
	return &db
}

func TestAPIMongoDB_AddRate(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(t, db)

	db.Init()
	if db.CountRate() != 0 {
		t.Error("database not properly initiated, should be empty")
	}
	r := Rate{}
	r.Base = "EUR"
	r.Date = "2017-10-31"
	r.Rates = make(map[string]float64)
	r.Rates["AUD"] = 1.5018
	r.Rates["BGN"] = 1.9558

	_ = db.AddRate(r)

	if db.CountRate() != 1 {
		t.Error("Did not insert properly")
	}
}

func TestAPIMongoDB_GetRate(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(t, db)
	db.Init()
	if db.CountRate() != 0 {
		t.Error("database not properly initiated, should be empty")
	}
	r := Rate{}
	r.Base = "EUR"
	r.Date = "2017-10-31"
	r.Rates = make(map[string]float64)
	r.Rates["AUD"] = 1.5018
	r.Rates["BGN"] = 1.9558
	_ = db.AddRate(r)
	if db.CountRate() != 1 {
		t.Error("Did not insert properly")
	}
	newRate, ok := db.GetRate(r.Date)
	if !ok {
		t.Error("Couldn't find entry")
	}
	if newRate.Base != r.Base || newRate.Date != r.Date {
		if len(newRate.Rates) != len(r.Rates) {
			t.Error("entries does not match")
		}
	}
}

func TestAPIMongoDB_DeleteRate(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(t, db)
	db.Init()
	if db.CountRate() != 0 {
		t.Error("database not properly initiated, should be empty")
	}
	r := Rate{}
	r.Base = "EUR"
	r.Date = "2017-10-31"
	r.Rates = make(map[string]float64)
	r.Rates["AUD"] = 1.5018
	r.Rates["BGN"] = 1.9558
	_ = db.AddRate(r)
	if db.CountRate() != 1 {
		t.Error("Did not insert properly")
	}
	newRate, ok := db.GetRate(r.Date)
	if !ok {
		t.Error("Couldn't find entry")
	}

	_ = db.DeleteRate(newRate)

	if db.CountRate() != 0 {
		t.Error("Did not delete entry")
	}
}

func TestGetRateFromAPI(t *testing.T) {
	rates := GetRateFromAPI()
	if rates.Date == "" {
		t.Error("GetRateFromAPI failed")
	} else {
		base := "EUR"
		target := "NOK"
		value, _ := ParseRate(base, target, rates)
		if rates.Rates[target] != value {
			t.Error("ParseRate for base EUR failed")
		}
		base = "NOK"
		target = "EUR"
		value, _ = ParseRate(base, target, rates)
		if value != 1/rates.Rates[base] {
			t.Error("ParseRate for target EUR failed")
		}
		target = "SEK"
		value, _ = ParseRate(base, target, rates)
		if value == 0 {
			t.Error("ParseRate for no EUR Probably didnt work")
		}
	}

}

func TestHandlerLatest(t *testing.T) {
	testStruct := TestInput{}
	testStruct.Result2.Parameters = make(map[string]string)
	testStruct.Result2.Parameters["baseCurrency"] = "EUR"
	testStruct.Result2.Parameters["targetCurrency"] = "NOK"

	body, err := json.Marshal(testStruct)
	if err != nil || body == nil {
		t.Error("marshalling failed")
	} else {
		log.Println(body)
		req, err := http.NewRequest("POST", "/latest/", bytes.NewBuffer(body))
		if err != nil {
			t.Error("nw request faield")
		} else {
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(HandlerLatest)
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != http.StatusOK {
				t.Error("handler not returning 200 ", status, http.StatusOK)
			}
		}
	}

}

//deleting the database after testing

func tearDownDB(t *testing.T, db *APIMongoDB) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		t.Error(err)
	}

	err = session.DB(db.DatabaseName).DropDatabase()
	if err != nil {
		t.Error(err)
	}
}
