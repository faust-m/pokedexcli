package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/faust-m/pokedexcli/internal/pokecache"
)

var cache *pokecache.Cache

func GetLocationAreas(requestURL string) (LocationArea, error) {
	if cache == nil {
		cache = pokecache.NewCache(5 * time.Second)
	}
	if result, found := cache.Get(requestURL); found {
		var locationData LocationArea
		err := json.Unmarshal(result, &locationData)
		if err != nil {
			return LocationArea{}, fmt.Errorf("error deserializing cached data: %w", err)
		}
		return locationData, nil
	}
	res, err := http.Get(requestURL)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error getting resource: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return LocationArea{}, fmt.Errorf("response returned with status: %s", res.Status)
	}
	decoder := json.NewDecoder(res.Body)
	var data LocationArea
	err = decoder.Decode(&data)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error decoding response: %w", err)
	}
	cacheData, err := json.Marshal(&data)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error serializing data: %w", err)
	}
	cache.Add(requestURL, cacheData)

	return data, nil
}
