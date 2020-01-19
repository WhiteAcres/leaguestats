package stats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/WhiteAcres/leaguestats/client"
	"github.com/WhiteAcres/leaguestats/storage"
)

type teamChampionPair struct {
	TeamID     int64
	ChampionID int64
}

type enemyChampionIDStatsObject struct {
	Name         string
	TotalMatches int64
	Victories    int64
	WinRate      float64
	BanScore     float64
}

type ddragonChampionPageObject struct {
	Kind    string `json:"type"`
	Format  string
	Version string
	Data    map[string]ddragonChampionObject
}

type ddragonChampionObject struct {
	ID   string
	Key  string
	Name string
}

func summonerInMatch(summonerName string, match client.Match) bool {
	participantsIdentities := match.ParticipantIdentities
	for _, participant := range participantsIdentities {
		player := participant.Player
		if player.SummonerName == summonerName {
			return true
		}
	}
	return false
}

func getParticipantIDForSummonerInMatch(summonerName string, match client.Match) int64 {
	participantsIdentities := match.ParticipantIdentities
	for _, participant := range participantsIdentities {
		player := participant.Player
		if player.SummonerName == summonerName {
			return participant.ParticipantID
		}
	}
	return -1
}

// returns true if the first gameVersion string is greater than the second, false otherwise
func versionCompare(gv1, gv2 string) bool {
	gv1Sections := strings.Split(gv1, ".")
	gv1S0, _ := strconv.ParseInt(gv1Sections[0], 10, 32)
	gv1S1, _ := strconv.ParseInt(gv1Sections[1], 10, 32)
	gv1S2, _ := strconv.ParseInt(gv1Sections[2], 10, 32)
	gv1S3, _ := strconv.ParseInt(gv1Sections[3], 10, 32)
	gv2Sections := strings.Split(gv2, ".")
	gv2S0, _ := strconv.ParseInt(gv2Sections[0], 10, 32)
	gv2S1, _ := strconv.ParseInt(gv2Sections[1], 10, 32)
	gv2S2, _ := strconv.ParseInt(gv2Sections[2], 10, 32)
	gv2S3, _ := strconv.ParseInt(gv2Sections[3], 10, 32)
	if gv1S0 > gv2S0 && gv1S1 > gv2S1 && gv1S2 > gv2S2 && gv1S3 > gv2S3 {
		return true
	}
	return false
}

// GetLatestGameVersion returns the latest gameVersion from all matches in storage
func GetLatestGameVersion(s storage.Storage) string {
	data := s.Data
	var gameVersion string = "0.0.0.0"
	for _, match := range data {
		if versionCompare(match.GameVersion, gameVersion) {
			gameVersion = match.GameVersion
		}
	}
	return gameVersion
}

// GetMatches gets all the matche in storage
func GetMatches(s storage.Storage) []client.Match {
	var matches []client.Match
	data := s.Data
	for _, match := range data {
		matches = append(matches, match)
	}
	return matches
}

// GetMatchesForSummoner gets all the matches for a summoner
func GetMatchesForSummoner(s storage.Storage, summonerName string) []client.Match {
	var summonerMatches []client.Match
	data := s.Data
	for _, match := range data {
		if summonerInMatch(summonerName, match) {
			summonerMatches = append(summonerMatches, match)
		}
	}
	return summonerMatches
}

// GetVictoryMatchesForSummoner gets all the victory matches for a summoner
func GetVictoryMatchesForSummoner(s storage.Storage, summonerName string) []client.Match {
	var summonerVictoryMatches []client.Match
	data := s.Data
	for _, match := range data {
		if summonerInMatch(summonerName, match) {
			summonerPID := getParticipantIDForSummonerInMatch(summonerName, match)
			victory := false
			for _, participant := range match.Participants {
				if participant.ParticipantID == summonerPID {
					victory = participant.Stats.Win
					break
				}
			}
			if victory {
				summonerVictoryMatches = append(summonerVictoryMatches, match)
			}
		}
	}
	return summonerVictoryMatches
}

// GetDefeatMatchesForSummoner gets all the defeat matches for a summoner
func GetDefeatMatchesForSummoner(s storage.Storage, summonerName string) []client.Match {
	var summonerDefeatMatches []client.Match
	data := s.Data
	for _, match := range data {
		if summonerInMatch(summonerName, match) {
			summonerPID := getParticipantIDForSummonerInMatch(summonerName, match)
			victory := false
			for _, participant := range match.Participants {
				if participant.ParticipantID == summonerPID {
					victory = participant.Stats.Win
					break
				}
			}
			if victory == false {
				summonerDefeatMatches = append(summonerDefeatMatches, match)
			}
		}
	}
	return summonerDefeatMatches
}

// GetEnemyChampionCountsInMatchesForSummoner returns count of all champions that summoner played against
func GetEnemyChampionCountsInMatchesForSummoner(summonerName string, matches []client.Match) map[int64]int64 {
	enemyChampCounts := make(map[int64]int64)
	for _, match := range matches {
		participantTeamMap := make(map[int64]*teamChampionPair)
		for _, participant := range match.Participants {
			participantTeamMap[participant.ParticipantID] = &teamChampionPair{participant.TeamID, participant.ChampionID}
		}
		summonerPID := getParticipantIDForSummonerInMatch(summonerName, match)
		summonerTID := participantTeamMap[summonerPID].TeamID
		for _, pair := range participantTeamMap {
			if pair.TeamID != summonerTID {
				champID := pair.ChampionID
				if val, ok := enemyChampCounts[champID]; ok {
					enemyChampCounts[champID] = val + 1
				} else {
					enemyChampCounts[champID] = 1
				}
			}
		}
	}
	return enemyChampCounts
}

// GetChampionCountsInMatches returns count of each champion from match list
func GetChampionCountsInMatches(matches []client.Match) map[int64]int64 {
	champCounts := make(map[int64]int64)
	for _, match := range matches {
		for _, participant := range match.Participants {
			champID := participant.ChampionID
			if val, ok := champCounts[champID]; ok {
				champCounts[champID] = val + 1
			} else {
				champCounts[champID] = 1
			}
		}
	}
	return champCounts
}

func getChampionWinRates(championCounts map[int64]int64, championVictoryCounts map[int64]int64) map[int64]float64 {
	championWinRates := make(map[int64]float64)
	for champID, totalCount := range championCounts {
		wins := float64(0)
		if val, ok := championVictoryCounts[champID]; ok {
			wins += float64(val)
		}
		winrate, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(wins)/float64(totalCount)), 64)
		championWinRates[champID] = winrate
	}
	return championWinRates
}

func calculateChampionBanScores(enemyChampCounts map[int64]int64, enemyChampWinRates map[int64]float64) map[int64]float64 {
	championBanScores := make(map[int64]float64)
	for champID, rate := range enemyChampWinRates {
		count := int64(1)
		if val, ok := enemyChampCounts[champID]; ok {
			count = val
		}
		adjRate := float64(rate) + float64(.5)
		score, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", math.Pow(adjRate, float64(count))), 64)
		championBanScores[champID] = score
	}
	return championBanScores
}

func getChampionNamesMap(latestGameVersion string) (map[int64]string, error) {
	// Creating the requess
	championNamesMap := make(map[int64]string)
	lgvSections := strings.Split(latestGameVersion, ".")
	lgvS0 := lgvSections[0]
	lgvS1 := lgvSections[1]
	urlstring := "http://ddragon.leagueoflegends.com/cdn/" + lgvS0 + "." + lgvS1 + ".1" + "/data/en_US/champion.json"
	fmt.Println(urlstring)
	req, err := http.NewRequest("GET", urlstring, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Translating response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var ddcpo ddragonChampionPageObject
	err = json.Unmarshal(body, &ddcpo)
	if err != nil {
		return nil, err
	}
	for _, v := range ddcpo.Data {
		champID, _ := strconv.ParseInt(v.Key, 10, 64)
		championNamesMap[champID] = v.Name
	}
	return championNamesMap, nil
}

// GetBestBanForSummoner does a thing
func GetBestBanForSummoner(s storage.Storage, summonerName string) {
	summonerMatches := GetMatchesForSummoner(s, summonerName)
	enemyChampCounts := GetEnemyChampionCountsInMatchesForSummoner(summonerName, summonerMatches)

	summonerDefeatMatches := GetDefeatMatchesForSummoner(s, summonerName)
	enemyChampVictoryCounts := GetEnemyChampionCountsInMatchesForSummoner(summonerName, summonerDefeatMatches)

	enemyChampWinRates := getChampionWinRates(enemyChampCounts, enemyChampVictoryCounts)
	enemyChampionBanScores := calculateChampionBanScores(enemyChampCounts, enemyChampWinRates)

	latestGameVersion := GetLatestGameVersion(s)
	championNamesMap, _ := getChampionNamesMap(latestGameVersion)

	var enemyChampionIDStatsList []enemyChampionIDStatsObject
	for champID, count := range enemyChampCounts {
		champName := "None"
		if val, ok := championNamesMap[champID]; ok {
			champName = val
		}
		victories := int64(0)
		if val, ok := enemyChampVictoryCounts[champID]; ok {
			victories = val
		}
		winrate := float64(0)
		if val, ok := enemyChampWinRates[champID]; ok {
			winrate = val
		}
		banscore := float64(0)
		if val, ok := enemyChampionBanScores[champID]; ok {
			banscore = val
		}
		eciso := enemyChampionIDStatsObject{champName, count, victories, winrate, banscore}
		enemyChampionIDStatsList = append(enemyChampionIDStatsList, eciso)
	}
	sort.Slice(enemyChampionIDStatsList, func(i, j int) bool {
		return enemyChampionIDStatsList[i].BanScore > enemyChampionIDStatsList[j].BanScore
	})
	fmt.Println(enemyChampionIDStatsList)
}
