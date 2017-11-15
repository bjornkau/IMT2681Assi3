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
//testing adding to the database
func TestAPIMongoDB_AddWebHook(t *testing.T){
	db := setUpDB(t)
	defer tearDownDB(t,db) 

	db.Init()
	
	if(db.CountWebHook() != 0){
		t.Error("database not properly initiated, should be empty")
	}

	w := WebHook{
		1,
		"example.com",
		"EUR",
		"NOK",
		2.655,
		3.68}
	
	_ = db.AddWebHook(w)

	if(db.CountWebHook() != 1 ){
		t.Error("Did not insert properly")
	}
}
//testing finding specific items in the database
func TestAPIMongoDB_GetWebHook(t *testing.T){
	db := setUpDB(t)
	defer tearDownDB(t,db) 

	db.Init()
	
	if(db.CountWebHook() != 0){
		t.Error("database not properly initiated, should be empty")
	}
	webHook := WebHook{
		1,
		"example.com",
		"EUR",
		"NOK",
		2.655,
		3.68}

	_ = db.AddWebHook(webHook)
	if db.CountWebHook() != 1 {
		t.Error("struct not added properly")
	}

	newWebhook, ok := db.GetWebHook(webHook.HookID)

	if (!ok) {
		t.Error("Couldn't find entry")
	}

	if (newWebhook.HookID != webHook.HookID ||
		newWebhook.WebHookURL != webHook.WebHookURL ||
		newWebhook.BaseCurrency != webHook.BaseCurrency ||
		newWebhook.TargetCurrency != webHook.TargetCurrency ||
		newWebhook.MinTriggerValue != webHook.MinTriggerValue ||
		newWebhook.MaxTriggerValue != webHook.MaxTriggerValue) {
			t.Error("does not match")
		}
}

//testing deletion from the database
func TestAPIMongoDB_DeleteWebHook(t *testing.T){
	db := setUpDB(t)
	defer tearDownDB(t,db) 

	db.Init()
	
	if(db.CountWebHook() != 0){
		t.Error("database not properly initiated, should be empty")
	}

	w := WebHook{
		1,
		"example.com",
		"EUR",
		"NOK",
		2.655,
		3.68}
	
	_ = db.AddWebHook(w)

	if(db.CountWebHook() != 1 ){
		t.Error("Did not insert properly")
	}

	newWebhook, ok := db.GetWebHook(w.HookID)

	if (!ok) {
		t.Error("Couldn't find entry")
	}

	_ = db.DeleteWebHook(newWebhook)

	if(db.CountWebHook() != 0){
		t.Error("Did not delete entry")
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
