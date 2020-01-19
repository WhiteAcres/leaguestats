package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/WhiteAcres/leaguestats/config"
	"github.com/WhiteAcres/leaguestats/stats"
	"github.com/WhiteAcres/leaguestats/storage"
)

func main() {

	conf := config.LoadConfig()
	conf.ValidateConfig()

	storage := storage.LoadStorage()

	// Get summoner name
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your Summoner Name:\n")
	summonerName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	summonerName = strings.Replace(summonerName, "\n", "", -1)
	summonerName = strings.Replace(summonerName, "\r", "", -1)

	// // initializing the client
	// url, _ := url.Parse("https://na1.api.riotgames.com")
	// cli := &client.Client{
	// 	BaseURL:    url,
	// 	APIKey:     conf.APIKey,
	// 	HTTPClient: &http.Client{}}

	// si, err := cli.GetSummonerInfo(summonerName)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Get the matches list
	// ml, err := cli.GetMatchList(si.AccountID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Pull out all gameIDs
	// var gameIDs []int64
	// for _, match := range ml.Matches {
	// 	// fmt.Println(match.GameID)
	// 	gameIDs = append(gameIDs, match.GameID)
	// }

	// // Get the match information for 50 the most recent gameIDs
	// var matches []*client.Match
	// for _, gameID := range gameIDs[0:50] {
	// 	m, err := cli.GetMatch(strconv.FormatInt(gameID, 10))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	// fmt.Println(m)
	// 	matches = append(matches, m)
	// }
	// storage.UpsertRecords(matches)

	stats.GetBestBanForSummoner(*storage, summonerName)
}
