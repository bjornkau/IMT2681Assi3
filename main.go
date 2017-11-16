package main

import(
	"io/ioutil"
	"encoding/json"
	"net/http"
	//"fmt"
	"strconv"
	"strings"
	"time"
	)

var ApiURL = string("http://api.fixer.io/latest")

var DISCLAIMER = string("As EUR is neither the base or target for currencies, a margin of error is to be expected.")

//Rate stores all information about one days rates
type Rate struct {
	Base string `json:"base"`
	Date string `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

type ResponsePayload struct {
	Speech string `json:"speech"`
	DisplayText string `json:"displayText"`
}

//post request Format "EUR to NOK"?
//HandlerLatest handles all querys to the bot
func HandlerLatest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		baseCurrency, targetCurrency := ParseInputBody(reqBody)
		value, baseEuro := GetResponse(baseCurrency, targetCurrency)
		respString := CreateResp(value, baseEuro, baseCurrency, targetCurrency)
		response := ResponsePayload{respString, respString}
		http.Header.Add(w.Header(), "content-type", "application/json")
		json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
		} else {

		}
	} else {
		http.Error(w, http.StatusText(400), 400)
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

func ParseInputBody(input map[string]interface{}) (baseCurrency string, targetCurrency string){
	result := input["result"].(map[string]interface{})
	parameters := result["parameters"].(map[string]interface{})
	baseCurrency = parameters["baseCurrency"].(string)
	targetCurrency = parameters["targetCurrency"].(string)
	return
}

func GetResponse(baseCurrency string, targetCurrency string) (value float64, baseEuro bool){
	db := SetUpDB()
	date := GetDate()
	rateFromDB, found := db.GetRate(date)
	if found {
		value, baseEuro = ParseRate(baseCurrency, targetCurrency, rateFromDB)
	} else {
		if(CheckTime()){
			rates := GetRateFromAPI()
			db.AddRate(rates)
			value, baseEuro = ParseRate(baseCurrency, targetCurrency, rates)
		} else {
			time := strings.Split(time.Now().AddDate(0,0,-1).String(), " ")
			DbRate, found1 := db.GetRate(time[0])
			if (found1) {
				value, baseEuro = ParseRate(baseCurrency, targetCurrency, DbRate)				
			} else {
				baseEuro = false
				value = -1.0
			}
		}
	}
	return
}

func ParseRate(baseCurrency string, targetCurrency string, rate Rate) (value float64, baseEuro bool) {
	baseEuro = false
	if baseCurrency == "EUR" {
		value = rate.Rates[targetCurrency]
		baseEuro = true
	} else if targetCurrency == "EUR" {
		euroToCurrency := rate.Rates[baseCurrency]
		value = 1/euroToCurrency
		baseEuro = true
	} else {
		euroToBase := rate.Rates[baseCurrency]
		euroToTarget := rate.Rates[targetCurrency]
		value = (1/euroToBase) * euroToTarget
	}
	return
}

func CreateResp(value float64, baseEuro bool, baseCurrency string, targetCurrency string) (respString string){
	respString = "The exchange rate between " + baseCurrency + " and " + targetCurrency + "is: "
	respString += strconv.FormatFloat(value, 'f', 4, 64) + "."
	if !baseEuro {
		respString += " " + DISCLAIMER
	}
	return
}

func main() {
	http.HandleFunc("/latest/", HandlerLatest)
	http.ListenAndServe("localhost:8080", nil)
}