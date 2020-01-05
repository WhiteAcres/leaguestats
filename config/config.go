package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"regexp"
	"strings"
)

// Conf - Config Object
type Conf struct {
	APIKey       string
	SummonerName string
}

// CreateDirIfNotExist is a thing
func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// getConfigFile returns a file pointer to conf file
func getConfigFile() *os.File {
	// Create the conf directory if not exists
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	confDirPath := user.HomeDir + "\\AppData\\Local\\leaguestats"
	createDirIfNotExist(confDirPath)

	// Read the conf file and create if it doesn't exist
	confFilePath := confDirPath + "\\conf.json"
	f, err := os.OpenFile(confFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// InitializeConfig initializes the config
func InitializeConfig() *Conf {
	f := getConfigFile()
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	var conf Conf
	err = json.Unmarshal(b, &conf)
	if err != nil || len(b) == 0 {
		return &Conf{}
	}
	return &conf
}

func validKey(apiKey string) bool {
	matched, err := regexp.MatchString(`RGAPI-\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`,
		apiKey)
	if err != nil {
		return false
	}
	return matched
}

func validSummonerName(name string) bool {
	return len(name) > 0
}

// ValidateConfig validates the config, updating the conf file if necessary
func (c *Conf) ValidateConfig() {
	if validKey(c.APIKey) == false {
		fmt.Println("Invalid API Key")
		fmt.Println("Generate new API Key at https://developer.riotgames.com")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the New API Key:\n")
		newAPIKey, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		newAPIKey = strings.Replace(newAPIKey, "\n", "", -1)
		newAPIKey = strings.Replace(newAPIKey, "\r", "", -1)
		c.APIKey = newAPIKey
	}
	if validSummonerName(c.SummonerName) == false {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Invalid Summoner Name")
		fmt.Print("Enter your Summoner Name:\n")
		summonerName, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		summonerName = strings.Replace(summonerName, "\n", "", -1)
		summonerName = strings.Replace(summonerName, "\r", "", -1)
		c.SummonerName = summonerName
	}
	fileData, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	f := getConfigFile()
	defer f.Close()
	_, err = f.Write(fileData)
	if err != nil {
		log.Fatal(err)
	}
}
