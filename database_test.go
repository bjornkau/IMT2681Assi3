package main

import (
"testing"
"gopkg.in/mgo.v2"
//"gopkg.in/mgo.v2/bson"
)

//setting up the database for testing
func setUpDB(t *testing.T) *APIMongoDB{
	db := APIMongoDB{
		"mongodb://localhost",
		"testAss2DB",
		"rates",
		"webHook",
	}

	session, err := mgo.Dial(db.Host)
	defer session.Close()

	if err != nil {
	t.Error(err)	
	}
	return &db
}
//deleting the database after testing
func tearDownDB(t *testing.T,db *APIMongoDB){
	session, err := mgo.Dial(db.Host)
	if err != nil {
	t.Error(err)		
	}

	err = session.DB(db.DatabaseName).DropDatabase()
	if(err != nil) {
		t.Error(err)
	}
}

func TestAPIMongoDB_AddRate(t *testing.T){
	db := setUpDB(t)
	defer tearDownDB(t, db)

	db.Init()
	if(db.CountRate() != 0){
		t.Error("database not properly initiated, should be empty")
	}
	r := Rate {}
	r.Base = "EUR"
	r.Date = "2017-10-31"
	r.Rates = make(map[string]float64)
	r.Rates["AUD"] = 1.5018
	r.Rates["BGN"] = 1.9558

	_ = db.AddRate(r)

	if (db.CountRate() != 1){
		t.Error("Did not insert properly")
	}
}

func TestAPIMongoDB_GetRate(t *testing.T){
	db := setUpDB(t)
	defer tearDownDB(t, db)
	db.Init()
	if(db.CountRate() != 0){
		t.Error("database not properly initiated, should be empty")
	}
	r := Rate {}
	r.Base = "EUR"
	r.Date = "2017-10-31"
	r.Rates = make(map[string]float64)
	r.Rates["AUD"] = 1.5018
	r.Rates["BGN"] = 1.9558
	_ = db.AddRate(r)
	if (db.CountRate() != 1){
		t.Error("Did not insert properly")
	}
	newRate, ok := db.GetRate(r.Date)
	if(!ok){
		t.Error("Couldn't find entry")
	}
	if(newRate.Base != r.Base || newRate.Date != r.Date){
		if (len(newRate.Rates) != len(r.Rates)){
			t.Error("entries does not match")
		}
	}
}

func TestAPIMongoDB_DeleteRate(t *testing.T){
	db := setUpDB(t)
	defer tearDownDB(t, db)
	db.Init()
	if(db.CountRate() != 0){
		t.Error("database not properly initiated, should be empty")
	}
	r := Rate {}
	r.Base = "EUR"
	r.Date = "2017-10-31"
	r.Rates = make(map[string]float64)
	r.Rates["AUD"] = 1.5018
	r.Rates["BGN"] = 1.9558
	_ = db.AddRate(r)
	if (db.CountRate() != 1){
		t.Error("Did not insert properly")
	}
	newRate, ok := db.GetRate(r.Date)
	if(!ok){
		t.Error("Couldn't find entry")
	}

	_ = db.DeleteRate(newRate)

	if(db.CountRate() != 0){
		t.Error("Did not delete entry")
	}
}
