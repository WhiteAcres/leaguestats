package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/WhiteAcres/leaguestats/client"
)

// Storage - json representation of all the Match objects
type Storage struct {
	Data map[int64]client.Match
}

// Need a load func (We may not actually need an InitilizeStorage func - we could just include empty init logic here)
func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getStorageFile() *os.File {
	// Create the conf directory if not exists
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	confDirPath := user.HomeDir + "\\AppData\\Local\\leaguestats"
	createDirIfNotExist(confDirPath)

	// Read the conf file and create if it doesn't exist
	confFilePath := confDirPath + "\\storage.json"
	f, err := os.OpenFile(confFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// LoadStorage loads the storage file
func LoadStorage() *Storage {
	f := getStorageFile()
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	var storage Storage
	err = json.Unmarshal(b, &storage)
	if err != nil || len(b) == 0 {
		return &Storage{make(map[int64]client.Match)}
	}
	return &storage
}

// SaveStorage saves the config file
func (s *Storage) SaveStorage() {
	fileData, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	f := getStorageFile()
	defer f.Close()
	_, err = f.Write(fileData)
	if err != nil {
		log.Fatal(err)
	}
}

// UpsertRecords inserts matches into the storage if they don't exist or updates them
func (s *Storage) UpsertRecords(matches []*client.Match) {
	for _, match := range matches {
		s.Data[match.GameID] = *match
	}
	s.SaveStorage()
}

// DeleteRecords deletes matches from the storage
func (s *Storage) DeleteRecords(matches []*client.Match) {
	for _, match := range matches {
		delete(s.Data, match.GameID)
	}
	s.SaveStorage()
}

// I might want to implement some sort of search functionality or maybe try to load the Storage object into a local MongoDB using a Mongo Library or something
// I don't think I really need this, though. I'll probably want to have the search logic be included in the actual operation function in "main.go"
// e.g. "func GetMostUsefulBan()" will implement logic to search through all matches for the current summonerName and see which champions are have the highest loss rate.
