package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

//APIMongoDB stores the details of the DB connection.
type APIMongoDB struct{
	Host string 
	DatabaseName string
	RatesCollectionName string
	WebHookCollectionName string
}

//Init checks if database exists
func (db *APIMongoDB) Init() {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
}

//AddRate takes param Rate and adds to the db
func (db *APIMongoDB) AddRate(r Rate) error {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	errInsert := session.DB(db.DatabaseName).C(db.RatesCollectionName).Insert(r)
	if errInsert != nil {
		fmt.Printf("Error in Insert(): %v", errInsert.Error())
		return errInsert
	}
	return nil
}

//CountRate returns amount of rates in DB
func (db *APIMongoDB) CountRate() int {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	count, errCount := session.DB(db.DatabaseName).C(db.RatesCollectionName).Count()
	if errCount != nil {
		fmt.Printf("Error in Count(): %v", errCount.Error())
		return -1
	}
	return count
}

//GetRate takes param Date (string) and returns Rate and bool
func (db *APIMongoDB) GetRate(Date string) (Rate, bool) {
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	rate := Rate{}
	allWasGood := true
	errFind := session.DB(db.DatabaseName).C(db.RatesCollectionName).Find(bson.M{"date": Date}).One(&rate)
	if errFind != nil {
		allWasGood = false
	}
	return rate, allWasGood
}

//DeleteRate takes param Rate and deletes from db, returns statusbool
func (db *APIMongoDB) DeleteRate(r Rate) (allIsWell bool){
	session, err := mgo.Dial(db.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	allIsWell = true
	err2 := session.DB(db.DatabaseName).C(db.RatesCollectionName).Remove(r)

	if err2 != nil{
		allIsWell = false
	}
	return
}