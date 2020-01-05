package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/WhiteAcres/leaguestats/client"
	"github.com/WhiteAcres/leaguestats/config"
)

func main() {

	conf := config.InitializeConfig()
	conf.ValidateConfig()

	// initializing the client
	url, _ := url.Parse("https://na1.api.riotgames.com")
	cli := &client.Client{
		BaseURL:    url,
		APIKey:     conf.APIKey,
		HTTPClient: &http.Client{}}

	si, err := cli.GetSummonerInfo(conf.SummonerName)
	if err != nil {
		log.Fatal(err)
	}

	// Get the matches list
	ml, err := cli.GetMatchList(si.AccountID)
	if err != nil {
		log.Fatal(err)
	}

	// Pull out all gameIDs
	var gameIDs []int64
	for _, match := range ml.Matches {
		fmt.Println(match.GameID)
		gameIDs = append(gameIDs, match.GameID)
	}

	// Get the match information for the most recent gameID
	mid := gameIDs[len(gameIDs)-1]
	m, err := cli.GetMatch(strconv.FormatInt(mid, 10))
	fmt.Println(m)
}
