package main

import(
	"io/ioutil"
	"encoding/json"
	"net/http"
	//"fmt"
	//"strconv"
	"strings"
	"time"
	)

var ApiURL = string("http://api.fixer.io/latest")

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
//post request Format "EUR to NOK"?
//HandlerLatest handles all querys to the bot
func HandlerLatest(w http.ResponseWriter, r *http.Request) {
	db := SetUpDB()
	date := GetDate()
	rateFromDB, found := db.GetRate(date)

	if(found){
		json.NewEncoder(w).Encode(rateFromDB.Rates["DDK"])
	} else {
		if(CheckTime()){
			rates := GetRateFromAPI()
			db.AddRate(rates)
			json.NewEncoder(w).Encode(rates.Rates["SEK"])
		} else {
			time := strings.Split(time.Now().AddDate(0,0,-1).String(), " ")
			DbRate, find := db.GetRate(time[0])
			if (find) {

				json.NewEncoder(w).Encode(DbRate.Rates["NOK"])
			}
		}
	}
}

//SetUpDB does set up db
func SetUpDB() *APIMongoDB{
	db := APIMongoDB{
		"mongodb://localhost",
		"testAssi3DB",
		"rates",
		"webHook",
	}
	return &db
}

//GetDate returns current date
func GetDate() (date string){
	timeStringParts := strings.Split(time.Now().String(), " ")
	date = timeStringParts[0]
	return
}

//Check time returs true if time is later than 1700CET
func CheckTime() (isLater bool){
	timeNow := time.Now().Hour()
	if(timeNow > 17){
		isLater = true
	} else {
		isLater = false
	}
	return
}

//GetRateFromAPI returns Rate with new info from api
func GetRateFromAPI() (rate Rate){
	response, err := http.Get(ApiURL)
	rate = Rate{}

	if(err != nil){

	}else {
		body,err := ioutil.ReadAll(response.Body)
		if err != nil{
			//TODO
		}else {
			err := json.Unmarshal(body, &rate)
			if err != nil{
				//TODO
			}

		}
	}
	return
}

func main() {
	http.HandleFunc("/latest/", HandlerLatest)
	http.ListenAndServe("localhost:8080", nil)
}