package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type Conf struct {
	APIKey       string
	SummonerName string
	BaseURL      string
}

type Client struct {
	BaseURL    *url.URL
	ApiKey     string
	httpClient *http.Client
}

type SummonerInfo struct {
	ID            string
	AccountID     string
	PuuID         string
	Name          string
	ProfileIconID int64
	RevisionDate  int64
	SummonerLevel int64
}

type MatchReference struct {
	Lane       string
	GameID     int64
	Champion   int64
	PlatformID string
	Season     int64
	Queue      int64
	Role       string
	Timestamp  int64
}

type Matchlist struct {
	Matches    []MatchReference `json:"matches"`
	TotalGames int64
	StartIndex int64
	EndIndex   int64
}

func (c *Client) getMatchList(accountID string) (*Matchlist, error) {
	// Creating the url
	rel := &url.URL{Path: "/lol/match/v4/matchlists/by-account/" + accountID}
	u := c.BaseURL.ResolveReference(rel)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("api_key", c.ApiKey)
	u.RawQuery = q.Encode()

	// Creating the request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	// Sending the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Translating response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var ml Matchlist
	err = json.Unmarshal(body, &ml)
	if err != nil {
		return nil, err
	}
	return &ml, nil
}

func (c *Client) getSummonerInfo(name string) (*SummonerInfo, error) {
	// Creating the url
	rel := &url.URL{Path: "/lol/summoner/v4/summoners/by-name/" + name}
	u := c.BaseURL.ResolveReference(rel)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("api_key", c.ApiKey)
	u.RawQuery = q.Encode()

	// Creating the request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	// Sending the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Translating response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var si SummonerInfo
	err = json.Unmarshal(body, &si)
	if err != nil {
		return nil, err
	}
	return &si, nil
}

func validKey(apiKey string) (bool, error) {
	matched, err := regexp.MatchString(`RGAPI-\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`,
		apiKey)
	if err != nil {
		return false, err
	}
	return matched, nil
}

func main() {
	// API Key available at https://developer.riotgames.com
	// reading in the config
	file, err := os.Open("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	var conf Conf
	err = json.Unmarshal(b, &conf)
	if err != nil {
		log.Fatal(err)
	}

	confUpdates := false

	// Validate the conf api key, potentially updating it
	vk, err := validKey(conf.APIKey)
	if err != nil {
		log.Fatal(err)
	}

	if vk == false {
		confUpdates = true
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
		conf.APIKey = newAPIKey
	}
	apikey := conf.APIKey

	// Validate the conf summoner name, potentially updating it
	if len(conf.SummonerName) == 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Invalid Summoner Name")
		fmt.Print("Enter your Summoner Name:\n")
		summonerName, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		summonerName = strings.Replace(summonerName, "\n", "", -1)
		summonerName = strings.Replace(summonerName, "\r", "", -1)
		conf.SummonerName = summonerName
	}

	if confUpdates {
		fileData, err := json.MarshalIndent(conf, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("conf.json", fileData, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// initializing the client
	url, _ := url.Parse(conf.BaseURL)
	client := &Client{
		BaseURL:    url,
		ApiKey:     apikey,
		httpClient: &http.Client{}}

	si, err := client.getSummonerInfo(conf.SummonerName)
	if err != nil {
		log.Fatal(err)
	}

	// Get the matches list
	ml, err := client.getMatchList(si.AccountID)
	if err != nil {
		log.Fatal(err)
	}

	// Pull out all gameIDs
	var gameIDs []int64
	for _, match := range ml.Matches {
		fmt.Println(match.GameID)
		gameIDs = append(gameIDs, match.GameID)
	}
}
