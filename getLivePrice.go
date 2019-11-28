package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var livePxCache map[string]float64 = map[string]float64{}

type Company struct {
	Name      string `json:"companyName"`
	Website   string `json:"website"`
	Employees int    `json:"employees"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
}

func autoPurgeLivePxCache() {
	for {
		time.Sleep(time.Minute * 5)

		livePxCache = map[string]float64{}
	}
}

func getLivePrice(symbol string) (float64, error) {
	cachedLast, ok := livePxCache[symbol]

	if ok {
		return cachedLast, nil
	}

	resp, err := http.Get("https://cloud.iexapis.com/v1/stock/" + symbol + "/quote/latestPrice?token=" + config["iexSecret"].(string))

	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return -1, err
	}

	floatPx, err := strconv.ParseFloat(string(data), 64)

	if err != nil {
		return -1, err
	}

	livePxCache[symbol] = floatPx

	return floatPx, nil
}

func getCompany(symbol string) (*Company, error) {
	resp, err := http.Get("https://cloud.iexapis.com/v1/stock/" + symbol + "/company?token=" + config["iexSecret"].(string))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	parsed := Company{}

	err = json.Unmarshal(data, &parsed)

	return &parsed, nil
}

func getLogo(symbol string) (string, error) {
	resp, err := http.Get("https://cloud.iexapis.com/v1/stock/" + symbol + "/logo?token=" + config["iexSecret"].(string))

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	parsed := struct {
		URL string `json:"url"`
	}{}

	err = json.Unmarshal(data, &parsed)

	return parsed.URL, nil
}

func getDayChange(symbol string) (float64, error) {
	resp, err := http.Get("https://cloud.iexapis.com/v1/stock/" + symbol + "/previous?token=" + config["iexSecret"].(string))

	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return -1, err
	}

	parsed := struct {
		Close float64 `json:"close"`
	}{}

	err = json.Unmarshal(data, &parsed)

	if err != nil {
		return -1, err
	}

	curPx, err := getLivePrice(symbol)

	if err != nil {
		return -1, err
	}

	changePercent := (curPx/parsed.Close - 1) * 100

	return changePercent, nil
}
