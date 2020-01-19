package stats

import (
	"fmt"
	"math"
	"strconv"

	"github.com/WhiteAcres/leaguestats/client"
	"github.com/WhiteAcres/leaguestats/storage"
)

type teamChampionPair struct {
	teamID     int64
	championID int64
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
		summonerTID := participantTeamMap[summonerPID].teamID
		for _, pair := range participantTeamMap {
			if pair.teamID != summonerTID {
				champID := pair.championID
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
		championWinRates[champID] = float64(wins) / float64(totalCount)
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
		score, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", math.Pow(adjRate, float64(count))), 64)
		championBanScores[champID] = score
	}
	return championBanScores
}

//TODO
func getChampionNamesMap() {
	// // Creating the request
	// req, err := http.NewRequest(method, u.String(), nil)
	// if err != nil {
	// 	return nil, err
	// }
	// client := http.Client
}

// GetBestBanForSummoner does a thing
func GetBestBanForSummoner(s storage.Storage, summonerName string) {
	summonerMatches := GetMatchesForSummoner(s, summonerName)
	enemyChampCounts := GetEnemyChampionCountsInMatchesForSummoner(summonerName, summonerMatches)

	summonerDefeatMatches := GetDefeatMatchesForSummoner(s, summonerName)
	enemyChampVictoryCounts := GetEnemyChampionCountsInMatchesForSummoner(summonerName, summonerDefeatMatches)

	enemyChampWinRates := getChampionWinRates(enemyChampCounts, enemyChampVictoryCounts)
	enemyChampionBanScores := calculateChampionBanScores(enemyChampCounts, enemyChampWinRates)
	fmt.Println(enemyChampionBanScores)
}
