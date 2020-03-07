package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/WhiteAcres/leaguestats/config"
)

// Client - League API Client Object
type Client struct {
	BaseURL    *url.URL
	APIKey     string
	HTTPClient *http.Client
}

// SummonerInfo - SummonerInfo Object from League API
type SummonerInfo struct {
	ID            string
	AccountID     string
	PuuID         string
	Name          string
	ProfileIconID int64
	RevisionDate  int64
	SummonerLevel int64
}

// MatchReference - MatchReference Object from League API
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

// Matchlist - Matchlist Object from League API
type Matchlist struct {
	Matches    []MatchReference
	TotalGames int64
	StartIndex int64
	EndIndex   int64
}

// Match - info regarding a partiular match
type Match struct {
	SeasonID              int64
	QueueID               int64
	GameID                int64
	ParticipantIdentities []ParticipantIdentity
	GameVersion           string
	PlatformID            string
	GameMode              string
	MapID                 int64
	GameType              string
	Teams                 []TeamStats
	Participants          []Participant
	GameDuration          int64
	GameCreation          int64
}

// ParticipantIdentity - info regarding paricipant identity
type ParticipantIdentity struct {
	Player        Player
	ParticipantID int64
}

// Player - info regarding player
type Player struct {
	CurrentPlatformID string
	SummonerName      string
	MatchHistoryURI   string
	PlatformID        string
	CurrentAccountID  string
	ProfileIcon       int
	SummonerID        string
	AccountID         string
}

// TeamStats - info regarding team stats
type TeamStats struct {
	FirstDragon          bool
	FirstInhibitor       bool
	Bans                 []TeamBans
	BaronKills           int64
	FirstRiftHerald      bool
	FirstBaron           bool
	RiftHeraldKills      int64
	FirstBlood           bool
	TeamID               int64
	FirstTower           bool
	VilemawKills         int64
	InhibitorKills       int64
	TowerKills           int64
	DominionVictoryScore int64
	Win                  string
	DragonKills          int64
}

// TeamBans - info regarding team bans
type TeamBans struct {
	PickTurn   int64
	ChampionID int64
}

// Participant - info regarding a participant
type Participant struct {
	Stats                     ParticipantStats
	ParticipantID             int64
	Runes                     []Rune
	Timeline                  ParticipantTimeline
	TeamID                    int64
	Spell2ID                  int64
	Masteries                 []Mastery
	HighestAchievedSeasonTier string
	Sepll1ID                  int64
	ChampionID                int64
}

// ParticipantStats - info regarding participant stats
type ParticipantStats struct {
	FirstBloodAssist                bool
	VisionScore                     int64
	MagicDamageDealtToChampions     int64
	DamageDealtToObjectives         int64
	TotalTimeCrowdControlDealt      int64
	LongestTimeSpentLiving          int64
	Perk1Var1                       int64
	Perk1Var3                       int64
	Perk1Var2                       int64
	TripleKills                     int64
	Perk3Var3                       int64
	NodeNeutralizeAssist            int64
	Perk3Var2                       int64
	PlayerScore9                    int64
	PlayerScore8                    int64
	kills                           int64
	PlayerScore1                    int64
	PlayerScore0                    int64
	PlayerScore3                    int64
	PlayerScore2                    int64
	PlayerScore5                    int64
	PlayerScore4                    int64
	PlayerScore7                    int64
	PlayerScore6                    int64
	Perk5Var1                       int64
	Perk5Var3                       int64
	Perk5Var2                       int64
	TotalScoreRank                  int64
	NeutralMinionsKilled            int64
	DamageDealtToTurrets            int64
	PhysicalDamageDealtToChampions  int64
	NodeCapture                     int64
	LargestMultiKill                int64
	Perk2Var2                       int64
	Perk2Var3                       int64
	TotalUnitsHealed                int64
	Perk2Var1                       int64
	Perk4Var1                       int64
	Perk4Var2                       int64
	Perk4Var3                       int64
	WardsKilled                     int64
	LargestCriticalStrike           int64
	LargestKillingSpree             int64
	QuadraKills                     int64
	TeamObjective                   int64
	MagicDamageDealt                int64
	Item2                           int64
	Item3                           int64
	Item0                           int64
	NeutralMinionsKilledTeamJungle  int64
	Item6                           int64
	Item4                           int64
	Item5                           int64
	Perk1                           int64
	Perk0                           int64
	Perk3                           int64
	Perk2                           int64
	Perk5                           int64
	Perk4                           int64
	Perk3Var1                       int64
	DamageSelfMitigated             int64
	MagicalDamageTaken              int64
	FirstInhibitorKilled            bool
	TrueDamageTaken                 int64
	NodeNeutralize                  int64
	Assists                         int64
	CombatPlayerScore               int64
	PerkPrimaryStyle                int64
	GoldSpent                       int64
	TrueDamageDealt                 int64
	ParticipantID                   int64
	TotalDamageTaken                int64
	PhysicalDamageDealt             int64
	SightWardsBoughtInGame          int64
	TotalDamageDealtToChampions     int64
	PhysicalDamageTaken             int64
	TotalPlayerScore                int64
	Win                             bool
	ObjectivePlayerScore            int64
	TotalDamageDealt                int64
	Item1                           int64
	NeutralMinionsKilledEnemyJungle int64
	Deaths                          int64
	WardsPlaced                     int64
	PerkSubStyle                    int64
	TurretKills                     int64
	FirstBloodKill                  bool
	TrueDamageDealtToChampions      int64
	GoldEarned                      int64
	KillingSprees                   int64
	UnrealKills                     int64
	AltersCaptured                  int64
	FirstTowerAssist                bool
	FirstTowerKill                  bool
	ChampLevel                      int64
	DoubleKills                     int64
	NodeCaptureAssist               int64
	InhibitorKills                  int64
	FirstInhibitorAssist            bool
	Perk0Var1                       int64
	Perk0Var2                       int64
	Perk0Var3                       int64
	VisionWardsBoughtInGame         int64
	AltarsNeutralized               int64
	PentaKills                      int64
	TotalHeal                       int64
	TotalMinionsKilled              int64
	TimeCCingOthers                 int64
}

// Rune - info regarding a rune
type Rune struct {
	RuneID int64
	Rank   int64
}

// ParticipantTimeline - info regarding a participant timeline
type ParticipantTimeline struct {
	Lane                        string
	ParticipantID               int64
	CSDiffPerMinuteDeltas       map[string]float64
	GoldPerMinDeltas            map[string]float64
	XPDiffPerMinDeltas          map[string]float64
	CreepsPerMinDeltas          map[string]float64
	XPPerMinDeltas              map[string]float64
	Role                        string
	DamageTakenDiffPerMinDeltas map[string]float64
	DamageTakenPerMinDeltas     map[string]float64
}

// Mastery - info regarding mastery
type Mastery struct {
	MasterID int64
	Rank     int64
}

// LeagueAPIRequest sends request to League API
func (c *Client) LeagueAPIRequest(method string, u *url.URL) ([]byte, error) {
	resp := &http.Response{}
	// Add api key to url
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("api_key", c.APIKey)
	u.RawQuery = q.Encode()

	// Creating the request
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	// Sending the request
	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	} else if resp.StatusCode == 403 {
		APIKey := config.GetNewAPIKey("API Key was unauthorized (probably expired)")
		c.APIKey = APIKey
		updates := map[string]string{"APIKey": APIKey}
		config.UpdateConfig(updates)
		time.Sleep(5 * time.Second)
	} else if resp.StatusCode == 404 {
		return nil, errors.New("Invalid Summoner Name")
	} else if resp.StatusCode == 429 {
		return nil, errors.New("API rate limit exceeded! Wait a few minutes, and try again.")
	} else if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return nil, errors.New("API Error")
	}

	// Translating response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return body, nil
}

// GetSummonerInfo - Gets Summoner Info from League API
func (c *Client) GetSummonerInfo(name string) (*SummonerInfo, error) {
	// Creating the url
	rel := &url.URL{Path: "/lol/summoner/v4/summoners/by-name/" + name}
	u := c.BaseURL.ResolveReference(rel)

	// Creating the request
	body, err := c.LeagueAPIRequest("GET", u)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var si SummonerInfo
	err = json.Unmarshal(body, &si)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &si, nil
}

// GetMatchList - Gets the match List from the League API
func (c *Client) GetMatchList(accountID string) (*Matchlist, error) {
	// Creating the url
	rel := &url.URL{Path: "/lol/match/v4/matchlists/by-account/" + accountID}
	u := c.BaseURL.ResolveReference(rel)

	body, err := c.LeagueAPIRequest("GET", u)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var ml Matchlist
	err = json.Unmarshal(body, &ml)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &ml, nil
}

// GetMatch - gets match information
func (c *Client) GetMatch(matchID string) (*Match, error) {
	// Creating the url
	rel := &url.URL{Path: "/lol/match/v4/matches/" + matchID}
	u := c.BaseURL.ResolveReference(rel)
	body, err := c.LeagueAPIRequest("GET", u)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var m Match
	err = json.Unmarshal(body, &m)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &m, nil
}
