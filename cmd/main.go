package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

const (
	//for the sake of simplicity i assume that there are no more than 20k teams starting from the same letter
	searchRequestTemplate   = "https://api.onefootball.com/entity-index-api/v1/archive/team/en-%s.json?page=1&limit=20000"
	teamDataRequestTemplate = "https://vintagemonster.onefootball.com/api/teams/en/%s.json"
)

var teamsList = []string{"Germany", "England", "France", "Spain", "Manchester United", "Arsenal", "Chelsea", "Barcelona", "Real Madrid", "Bayern Munich"}

//TeamData describes structure of team data response
type TeamData struct {
	Code int `json:"code"`
	Data struct {
		Team Team `json:"team"`
	}
}

//Team represents team entity
type Team struct {
	ID      json.Number `json:"id"`
	Name    string      `json:"name"`
	Players []Player    `json:"players"`
}

//Player represents player entity
type Player struct {
	ID    json.Number `json:"id"`
	Name  string      `json:"name"`
	Age   json.Number `json:"age"`
	Teams []string    `json:",omitempty"`
}

//Search is a struct for a search query response
type Search struct {
	Teams []struct {
		ID   json.Number
		Name string
	} `json:"data"`
}

func main() {
	teams := findTeams(teamsList)

	players, playersTeams := getPlayersListWithTeams(teams)

	for i := range players {
		fmt.Printf("%d. %s; %s; %s \n", i+1, players[i].Name, players[i].Age, strings.Join(playersTeams[players[i].ID.String()], ", "))
	}
}

func getPlayersListWithTeams(teams map[string]*Team) ([]Player, map[string][]string) {
	players := make([]Player, 0, 20*len(teams)) //capacity is just an assumption. football team consists of at least 20 players, AFAIK
	playersTeams := make(map[string][]string)

	for _, t := range teams {
		for _, p := range t.Players {
			if _, ok := playersTeams[p.ID.String()]; ok {
				playersTeams[p.ID.String()] = append(playersTeams[p.ID.String()], t.Name)
				continue
			}

			players = append(players, p)
			playersTeams[p.ID.String()] = append(playersTeams[p.ID.String()], t.Name)
		}
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].Name < players[j].Name
	})

	return players, playersTeams
}

func findTeams(list []string) map[string]*Team {
	teamsData := make(map[string]*Team)
	uniqueLetters := make(map[string]struct{})

	letters := make([]string, 0)
	for i := range list {
		l := list[i][:1]
		teamsData[list[i]] = nil
		if _, ok := uniqueLetters[l]; ok {
			continue
		}
		letters = append(letters, l)
		uniqueLetters[l] = struct{}{}
	}

	for i := range letters {
		res, err := httpCall(fmt.Sprintf(searchRequestTemplate, strings.ToLower(letters[i])))
		if err != nil {
			log.Println("error getting search data on ", letters[i], err)
			continue
		}

		var search Search
		if err := json.Unmarshal(res, &search); err != nil {
			log.Println("error unmarshal response on ", letters[i], err)
			continue
		}

		for _, t := range search.Teams {
			if _, ok := teamsData[t.Name]; !ok {
				continue
			}

			teamsData[t.Name], err = getTeamData(t.ID.String())
			if err != nil {
				log.Println("failed to get data for ", t.Name, err)
				continue
			}
		}
	}

	return teamsData
}

func getTeamData(id string) (*Team, error) {
	res, err := httpCall(fmt.Sprintf(teamDataRequestTemplate, id))
	if err != nil {
		return nil, err
	}

	var teamData TeamData
	if err := json.Unmarshal(res, &teamData); err != nil {
		return nil, err
	}

	return &teamData.Data.Team, nil
}

func httpCall(url string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}
