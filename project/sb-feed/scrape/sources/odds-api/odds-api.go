package oddsapisrc

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/MartinNikolovMarinov/sb-infra/entities"
	"github.com/tidwall/gjson"
)

const (
	// ConversionError signals data conversion errors
	ConversionError = "Conversion Error"
)

// OddsAPIFeed is the odds feed for https://app.oddsapi.io/api/v1/odds
type OddsAPIFeed struct {
	client    *http.Client
	SourceURL string
}

// New creates an instance of the OddsApiFeed struct
func New(client *http.Client, sourceURL string) *OddsAPIFeed {
	return &OddsAPIFeed{
		client:    client,
		SourceURL: sourceURL,
	}
}

// Fetch comes from the OddsFeed interface
func (feed *OddsAPIFeed) Fetch() ([]*entities.Game, error) {
	var resp, resp2 *http.Response
	var err error
	var games []*entities.Game
	var leagueToCountry map[string]string
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		resp, err = feed.client.Get(feed.SourceURL + "all-soccer.json")
		if err != nil {
			return
		}
		defer resp.Body.Close()
		games, err = feed.mapOddsResponse(resp.Body)
		if err != nil {
			return
		}
	}()

	go func() {
		defer wg.Done()
		resp2, err = feed.client.Get(feed.SourceURL + "leagues.json")
		if err != nil {
			return
		}
		defer resp2.Body.Close()

		leagueToCountry, err = feed.mapLeagueResponse(resp2.Body)
		if err != nil {
			return
		}
	}()

	wg.Wait()

	if err != nil {
		return nil, err
	}

	for _, g := range games {
		if country, ok := leagueToCountry[g.League.Name]; ok {
			g.League.Country = country
		}
	}

	return games, nil
}

/*** SCRAPPER CODE ***/

func (feed *OddsAPIFeed) mapLeagueResponse(reader io.Reader) (map[string]string, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var leaguesRaw []struct {
		Name    string
		Sport   string
		Country string
	}
	err = json.Unmarshal(bytes, &leaguesRaw)
	if err != nil {
		return nil, err
	}

	leagueToCountry := make(map[string]string)
	for _, l := range leaguesRaw {
		if l.Country != "" {
			leagueToCountry[l.Name] = l.Country
		}
	}

	return leagueToCountry, nil
}

func (feed *OddsAPIFeed) mapOddsResponse(response io.Reader) ([]*entities.Game, error) {
	bytes, err := ioutil.ReadAll(response)
	if err != nil {
		return nil, err
	}

	parsed, err := parseResponse(bytes)
	if err != nil {
		return nil, err
	}

	games := make([]*entities.Game, 0)
	for _, p := range parsed {
		sport := entities.NewSport(-1, p.Sport.Name)
		league := entities.NewLeague(-1, p.League.Name, sport, "")
		game := entities.NewGame(-1, p.Event.Home+" vs "+p.Event.Away, sport, league)
		game.AddTeam(entities.NewTeam(-1, p.Event.Home, league))
		game.AddTeam(entities.NewTeam(-1, p.Event.Away, league))
		ge := entities.NewGameEvent(-1, p.Event.StartDate, game, entities.T1x2)
		l := entities.NewLine(-1, entities.MoneyLine, "1x2", ge)

		for _, s := range p.Sites {
			odd := entities.NewOdd(-1, s.Name, s.Odds, l)
			l.AddOdd(odd)
		}

		ge.AddLine(l)
		game.AddEvent(ge)
		games = append(games, game)
	}

	return games, nil
}

// OddsAPIResponseObject the response object from https://app.oddsapi.io/api/v1/odds
type oddsAPIResponseObject struct {
	Sport struct {
		Key  string
		Name string
	}
	League struct {
		Key  string
		Name string
	}
	Event       oddsAPIEvent
	LastUpdated time.Time `json:"last_updated"`
	Sites       []oddsAPISite
}

type oddsAPIEvent struct {
	Away      string
	Home      string
	StartDate time.Time `json:"start_time"`
}

type oddsAPISite struct {
	Odds       []float32
	Bookmaker  bool
	Exchange   bool
	Name       string
	IsOutright bool
}

func parseResponse(bytes []byte) ([]oddsAPIResponseObject, error) {
	var err error
	rawData := gjson.ParseBytes(bytes).Array()
	result := make([]oddsAPIResponseObject, 0)

	for _, rawObjet := range rawData {
		object := oddsAPIResponseObject{}
		object.Sites = make([]oddsAPISite, 0)

		rawObjet.ForEach(func(key, value gjson.Result) bool {
			b, err2 := parseRoot(key, value, &object)
			if err2 != nil {
				err = err2
			}
			return b
		})

		result = append(result, object)
	}

	return result, err
}

func parseRoot(key, value gjson.Result, object *oddsAPIResponseObject) (bool, error) {
	var err error

	switch {
	case key.Type == gjson.String && key.Str == "sport":
		err = json.Unmarshal([]byte(value.Raw), &object.Sport)
	case key.Type == gjson.String && key.Str == "league":
		err = json.Unmarshal([]byte(value.Raw), &object.League)
	case key.Type == gjson.String && key.Str == "event":
		err = json.Unmarshal([]byte(value.Raw), &object.Event)
	case key.Type == gjson.String && key.Str == "sites":
		value.ForEach(func(key, value gjson.Result) bool {
			b, err2 := parseSiteEvents(key, value, object)
			if err2 != nil {
				err = err2
			}
			return b
		})
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func parseSiteEvents(key, value gjson.Result, object *oddsAPIResponseObject) (bool, error) {
	var err error

	switch {
	case key.Type == gjson.String && key.Str == "last_updated":
		err = json.Unmarshal([]byte(value.Raw), &object.LastUpdated)
	case key.Type == gjson.String && key.Str == "1x2":
		isOutright := false
		value.ForEach(func(key, value gjson.Result) bool {
			b, err2 := parseSites(key, value, object, &isOutright)
			if err2 != nil {
				err = err2
			}
			return b
		})
	default:
		// NOTE: skipping all other even types
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func parseSites(key, value gjson.Result, object *oddsAPIResponseObject, isOutright *bool) (bool, error) {
	var err error

	if key.Type == gjson.String && key.Str == "outright" {
		var ok bool
		*isOutright, ok = value.Value().(bool)
		if !ok {
			return false, errors.New(ConversionError)
		}
		return true, nil
	}

	sites := struct {
		Odds       []float32
		Bookmaker  bool
		Exchange   bool
		Name       string
		IsOutright bool
	}{}
	sites.Odds = make([]float32, 0)
	sites.IsOutright = *isOutright

	value.ForEach(func(key, value gjson.Result) bool {
		switch {
		case key.Type == gjson.String && key.Str == "odds":
			oddsMap, ok := value.Value().(map[string]interface{})
			if !ok {
				err = errors.New(ConversionError)
				return false
			}

			for _, odd := range oddsMap {
				asFloat64, ok := odd.(float64) // gjson converts js numbers to float64
				if !ok {
					err = errors.New(ConversionError)
					return false
				}
				sites.Odds = append(sites.Odds, float32(asFloat64))
			}
		case key.Type == gjson.String && key.Str == "bookmaker":
			asBool, ok := value.Value().(bool)
			if !ok {
				err = errors.New(ConversionError)
				return false
			}
			sites.Bookmaker = asBool
		case key.Type == gjson.String && key.Str == "exchange":
			asBool, ok := value.Value().(bool)
			if !ok {
				err = errors.New(ConversionError)
				return false
			}
			sites.Exchange = asBool
		case key.Type == gjson.String && key.Str == "name":
			asString, ok := value.Value().(string)
			if !ok {
				err = errors.New(ConversionError)
				return false
			}
			sites.Name = asString
		}

		return true
	})

	object.Sites = append(object.Sites, sites)

	if err != nil {
		return false, err
	}

	return true, nil
}
