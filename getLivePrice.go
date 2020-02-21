package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var livePxCache map[string]float64 = map[string]float64{}
var previousCache map[string]float64 = map[string]float64{}
var companyCache map[string]*Company = map[string]*Company{}

// Company is a company
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

func autoPurgePreviousCache() {
	for {
		time.Sleep(time.Minute * 5)

		previousCache = map[string]float64{}
	}
}

func autoPurgeCompanyCache() {
	for {
		time.Sleep(time.Hour * 24)

		companyCache = map[string]*Company{}
	}
}

func isCrypto(symbol string) bool {
	return strings.HasPrefix(symbol, "C:")
}

func getLastUpdated(symbol string) (*time.Time, error) {
	resp, err := http.Get("https://cloud.iexapis.com/v1/stock/" + strings.ToLower(symbol) + "/quote/latestUpdate?token=" + config["iexSecret"].(string))

	if err != nil {
		return nil, err
	}

	rd, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	prd, err := strconv.Atoi(string(rd))

	if err != nil {
		return nil, err
	}

	lastUpdate := time.Unix(0, int64(time.Millisecond)*int64(prd))

	return &lastUpdate, nil
}

func getLivePrice(symbol string, bypassCache bool) (float64, error) {
	if bypassCache == false {
		cachedLast, ok := livePxCache[symbol]

		if ok {
			return cachedLast, nil
		}
	}

	floatPx := 0.0

	parseAsCrypto := false

	var resp *http.Response
	var err error

	if isCrypto(symbol) {
		// Fetch crypto price

		resp, err = http.Get("https://cloud.iexapis.com/v1/crypto/" + strings.ToLower(string([]rune(symbol)[2:])) + "usd/price?token=" + config["iexSecret"].(string))

		parseAsCrypto = true
	} else {
		// Fetch stock price

		resp, err = http.Get("https://cloud.iexapis.com/v1/stock/" + strings.ToLower(symbol) + "/quote/latestPrice?token=" + config["iexSecret"].(string))
	}

	if err != nil {
		return -1, err
	}

	if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case 404:
			return -1, errors.New("The symbol does not exist.")
		default:
			return -1, errors.New("Server error " + strconv.Itoa(resp.StatusCode))
		}
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return -1, err
	}

	if parseAsCrypto {
		// Parse crypto price JSON

		parsed := struct {
			Price  string `json:"price"`
			Symbol string `json:"symbol"`
		}{}

		err = json.Unmarshal(data, &parsed)

		if err != nil {
			return -1, err
		}

		floatPx, err = strconv.ParseFloat(parsed.Price, 64)

		if err != nil {
			return -1, err
		}
	} else {
		// Parse float stock quote live price

		floatPx, err = strconv.ParseFloat(string(data), 64)

		if err != nil {
			return -1, err
		}
	}

	livePxCache[symbol] = floatPx

	return floatPx, nil
}

func getCompany(symbol string) (*Company, error) {
	if isCrypto(symbol) {
		// Not implemented for crypto

		return &Company{
			Name:      "Crypto: " + strings.ToUpper(string([]rune(symbol)[2:])),
			Website:   "N/A",
			Employees: 0,
			City:      "N/A",
			State:     "N/A",
			Country:   "BITCONNECT",
		}, nil
	}

	cachedCompany, ok := companyCache[symbol]

	if ok {
		return cachedCompany, nil
	}

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

	if err != nil {
		return nil, err
	}

	companyCache[symbol] = &parsed

	return &parsed, nil
}

func getLogo(symbol string) (string, error) {
	if isCrypto(symbol) {
		// Not implemented for crypto

		return "", nil
	}

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

	if err != nil {
		return "", err
	}

	return parsed.URL, nil
}

func getDayChange(symbol string) (float64, error) {
	if isCrypto(symbol) {
		// Day change not yet implemented for crypto

		return 0.0, nil
	}

	prevPx, ok := previousCache[symbol]

	if !ok {
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

		prevPx = parsed.Close

		previousCache[symbol] = prevPx
	}

	curPx, err := getLivePrice(symbol, false)

	if err != nil {
		return -1, err
	}

	changePercent := (curPx/prevPx - 1) * 100

	return changePercent, nil
}
