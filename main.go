package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	HTTP_PORT = ":5555"
)

type CurrentWeather struct {
	Condition string       `json:"condition"`
	Type      string       `json:"temperature"`
	Alerts    []AlertsType `json:"alerts,omitempty"`
}

type TemperatureType struct {
	Temp float32 `json:"temp"`
}

type ConditionType struct {
	Id          int32  `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type CurrentType struct {
	Temp      float32         `json:"temp"`
	Condition []ConditionType `json:"weather"`
}

type AlertsType struct {
	SenderName  string `json:"sender_name"`
	Event       string `json:"event"`
	Start       int64  `json:"start"`
	End         int64  `json:"end"`
	Description string `json:"description"`
}
type CurrentWeatherInfo struct {
	Current CurrentType  `json:"current"`
	Alerts  []AlertsType `json:"alerts"`
}

func main() {

	var wg sync.WaitGroup

	wg.Add(1)

	go startServer(&wg)

	wg.Wait()

}

func startServer(wg *sync.WaitGroup) {
	defer wg.Done()
	http.HandleFunc("/myApp/v1/weather/current", getCurrentWeather)
	log.Fatal(http.ListenAndServe(HTTP_PORT, nil))
}

func getCurrentWeather(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting current weather")
	lat, ok := r.URL.Query()["lat"]
	if !ok {
		log.Println("lat is a required parameter")
		w.WriteHeader(412)
		fmt.Fprintln(w, "required parameter is missing")
		return
	}
	long, ok := r.URL.Query()["long"]
	if !ok {
		log.Println("long is a required parameter")
		w.WriteHeader(412)
		fmt.Fprintln(w, "required parameter is missing")
		return
	}
	log.Println("lat:", lat[0], "long:", long[0])

	// get weather info from weather service
	// build http request
	url := "https://api.openweathermap.org/data/2.5/onecall"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	q := req.URL.Query()
	q.Add("lat", lat[0])
	q.Add("lon", long[0])
	q.Add("appid", "78fe259766f8b3433838fb022bf5724a")
	q.Add("units", "imperial")
	q.Add("exclude", "minutely,hourly,daily")

	req.URL.RawQuery = q.Encode()
	log.Println(req.URL.String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	if resp.StatusCode != 200 {
		w.WriteHeader(resp.StatusCode)
		fmt.Fprintln(w, resp.Status)
		return

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(resp.StatusCode)
		fmt.Fprintln(w, "Error reading response")
		return
	}
	log.Println(string(body))

	var currWeather CurrentWeatherInfo

	if err = json.Unmarshal(body, &currWeather); err != nil {
		log.Println("Error unmarshaling", err)
		return
	}

	jsonStr, err := json.Marshal(currWeather)
	if err != nil {
		log.Println("Error marshaling", err)
		return
	}
	log.Println(string(jsonStr))
	var currWeatherResult CurrentWeather
	// condition rain, snow, ...
	if len(currWeather.Current.Condition) > 0 {
		currWeatherResult.Condition = strings.ToLower(currWeather.Current.Condition[0].Main)
	}
	tempType := getTempType(currWeather.Current.Temp)
	currWeatherResult.Type = tempType
	if len(currWeather.Alerts) >= 0 {
		currWeatherResult.Alerts = make([]AlertsType, len(currWeather.Alerts))
		copy(currWeatherResult.Alerts, currWeather.Alerts)
	}

	jsonStrResult, err := json.Marshal(currWeatherResult)
	if err != nil {
		log.Println("Error marshaling", err)
		return
	}
	log.Println("Returning result....")
	log.Println(string(jsonStrResult))
	fmt.Fprintf(w, string(jsonStrResult)+"\n")
}

// getTempType get temperature type based on current temperature
func getTempType(temp float32) string {
	var tempType string
	if temp < 65.0 {
		tempType = "cold"
	} else if temp > 80 {
		tempType = "hot"
	} else {
		tempType = "moderate"
	}
	return tempType
}
