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

// // Prune modifies storage to only keep the most recent total number of records
// func (s *Storage) Prune(total int) {

// 	// Dump all stored matches into a slice
// 	var matches []*client.Match
// 	for _, match := range s.Data {
// 		matches = append(matches, &match)
// 	}

// 	// Sort all the matches (with the most recent being at the lowest slice index)
// 	sort.Slice(matches, func(i, j int) bool {
// 		return matches[i].GameCreation > matches[j].GameCreation
// 	})

// 	//Only keep the latest total number of games
// 	fmt.Println(len(matches))
// 	if len(matches) > total {
// 		matches = matches[0:total]
// 		newS := Storage{make(map[int64]client.Match)}
// 		for _, match := range matches {
// 			newS.Data[match.GameID] = *match
// 		}

// 		// Overwrite the old storage with the new one
// 		newS.SaveStorage()
// 	}
// }

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

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// FilterGameIDs returns a slice of gameIDs not already found in storage
func (s *Storage) FilterGameIDs(gameIDs []int64) []int64 {
	var filteredGameIDs []int64
	sGameIDs := make([]int64, 0, len(s.Data))
	for sGameID := range s.Data {
		sGameIDs = append(sGameIDs, sGameID)
	}
	for _, gameID := range gameIDs {
		if contains(sGameIDs, gameID) == false {
			filteredGameIDs = append(filteredGameIDs, gameID)
		}
	}
	return filteredGameIDs
}
