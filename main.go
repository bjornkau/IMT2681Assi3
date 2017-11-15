package main

import(
	"encoding/json"
	"net/http"
	"fmt"
	"strconv"
	"strings"
	"time"
	)



//WebHook stores all information about webhooks
type WebHook struct {
	HookID int `json:"hookId"`
	WebHookURL string `json:"webhookURL"`
	BaseCurrency string `json:"baseCurrency"`
	TargetCurrency string `json:"targetCurrency"`
	MinTriggerValue float64 `json:"minTriggerValue"`
	MaxTriggerValue float64 `json:"maxTriggerValue"`
}

//Rate stores all information about one days rates
type Rate struct {
	Base string `json:"base"`
	Date string `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

//HandlerLatest handles all querys tot he bot
func HandlerLatest(w http.ResponseWriter, r *http.Request) {
	db := SetUpDB()
	timeNow := time.Now()
	
	
}

func SetUpDB() *APIMongoDB{
	db := APIMongoDB{
		"mongodb://localhost",
		"testAssi3DB",
		"rates",
		"webHook",
	}
	return &db
}

func main() {
	http.HandleFunc("/latest/", HandlerLatest)
	http.ListenAndServe("localhost:8080", nil)
}