package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// StopOrder data
type StopOrder struct {
	Count int     `json:"count"`
	Price float64 `json:"price"`
}

// User data
type User struct {
	mux           sync.Mutex
	ID            string               `json:"id"`
	Balance       float64              `json:"balance"`
	Shares        map[string]int       `json:"shares"`
	LastReset     time.Time            `json:"lastReset"`
	Stops         map[string]StopOrder `json:"stops"`
	SeasonStartNW float64              `json:"seasonStartNetWorth"`
}

var memUsers map[string]*User = map[string]*User{}

func getUserFilePath(id string) (*string, error) {
	userFilePath, err := filepath.Abs(filepath.Join("users", id+".json"))

	if err != nil {
		return nil, err
	}

	return &userFilePath, nil
}

// ListUsers returns a list of user IDs
func ListUsers() ([]string, error) {
	usersDir, err := filepath.Abs("users")

	if err != nil {
		return nil, err
	}

	d, err := os.Open(usersDir)

	if err != nil {
		return nil, err
	}

	files, err := d.Readdir(-1)

	if err != nil {
		return nil, err
	}

	ids := []string{}

	for _, f := range files {
		fname := []rune(f.Name())
		ids = append(ids, string(fname[0:len(fname)-5]))
	}

	return ids, nil
}

// GetUser reads a User from its ID
func GetUser(id string) (*User, error) {
	memUser, ok := memUsers[id]

	var user User

	if ok {
		return memUser, nil
	}

	userFilePath, err := getUserFilePath(id)

	if err != nil {
		return nil, err
	}

	userFile, err := os.Open(*userFilePath)

	defer userFile.Close()

	if err != nil {
		// Create user if not existent

		user = User{
			ID:        id,
			mux:       sync.Mutex{},
			Balance:   config["game"].(map[string]interface{})["startBalance"].(float64),
			Shares:    map[string]int{},
			LastReset: time.Now(),
			Stops:     map[string]StopOrder{},
		}

		user.SeasonStartNW = user.Balance

		err := user.Save()

		memUsers[id] = &user

		if err != nil {
			return &user, nil
		}

		return nil, err
	}

	userFileContent, err := ioutil.ReadAll(userFile)

	if err != nil {
		return nil, err
	}

	user = User{
		mux: sync.Mutex{},
	}

	err = json.Unmarshal(userFileContent, &user)

	if err != nil {
		return nil, err
	}

	// Ensure user has Stops map

	if user.Stops == nil {
		user.Stops = map[string]StopOrder{}
	}

	// Update user

	err = UpdateUser(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Save user data
func (u *User) Save() error {
	userFilePath, err := getUserFilePath(u.ID)

	if err != nil {
		return err
	}

	file, err := os.Create(*userFilePath)
	defer file.Close()

	if err != nil {
		return err
	}

	marshalled, err := json.Marshal(u)

	if err != nil {
		return err
	}

	_, err = file.Write(marshalled)

	if err != nil {
		return err
	}

	return nil
}

// NetWorth returns User u's net worth
func (u *User) NetWorth(bypassCache bool) (float64, error) {
	nw := u.Balance

	for symbol, count := range u.Shares {
		px, err := getLivePrice(symbol, bypassCache)

		if err != nil {
			return 0, err
		}

		nw += px * float64(count)
	}

	return nw, nil
}

// UpdateUser completes all pending transactions for u
func UpdateUser(u *User) error {
	// Fulfill stop orders

	for symbol, stop := range u.Stops {
		stopSymbolPx, err := getLivePrice(symbol, false)

		if err != nil {
			return err
		}

		if stopSymbolPx < stop.Price {
			// Price has fallen below stop price
			// Execute stop at saved price due to IVT

			fmt.Println("Executing stop order", symbol, stop.Count, "at price", stop.Price)

			u.Balance += stop.Price * float64(stop.Count)
			u.Shares[symbol] -= stop.Count

			// Drop from Shares if count is 0

			if u.Shares[symbol] == 0 {
				delete(u.Shares, symbol)
			}

			delete(u.Stops, symbol)
		}
	}

	err := u.Save()

	if err != nil {
		return err
	}

	return nil
}
