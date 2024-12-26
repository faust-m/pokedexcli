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
		cache = pokecache.NewCache(5 * time.Minute)
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

func ExploreArea(requestURL string) (ExploreResult, error) {
	if cache == nil {
		cache = pokecache.NewCache(5 * time.Minute)
	}
	if result, found := cache.Get(requestURL); found {
		var exploreData ExploreResult
		err := json.Unmarshal(result, &exploreData)
		if err != nil {
			return ExploreResult{}, fmt.Errorf("error deserializing cached data: %w", err)
		}
		return exploreData, nil
	}
	res, err := http.Get(requestURL)
	if err != nil {
		return ExploreResult{}, fmt.Errorf("error getting resource: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ExploreResult{}, fmt.Errorf("response returned with status: %s", res.Status)
	}
	decoder := json.NewDecoder(res.Body)
	var data ExploreResult
	err = decoder.Decode(&data)
	if err != nil {
		return ExploreResult{}, fmt.Errorf("error decoding response: %w", err)
	}
	cacheData, err := json.Marshal(&data)
	if err != nil {
		return ExploreResult{}, fmt.Errorf("error serializing data: %w", err)
	}
	cache.Add(requestURL, cacheData)

	return data, nil
}

func GetPokemonData(requestURL string) (Pokemon, error) {
	if cache == nil {
		cache = pokecache.NewCache(5 * time.Minute)
	}
	if result, found := cache.Get(requestURL); found {
		var pokemonData Pokemon
		err := json.Unmarshal(result, &pokemonData)
		if err != nil {
			return Pokemon{}, fmt.Errorf("error deserializing cached data: %w", err)
		}
		return pokemonData, nil
	}
	res, err := http.Get(requestURL)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error getting resource: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Pokemon{}, fmt.Errorf("response returned with status: %s", res.Status)
	}
	decoder := json.NewDecoder(res.Body)
	var data Pokemon
	err = decoder.Decode(&data)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error decoding response: %w", err)
	}
	cacheData, err := json.Marshal(&data)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error serializing data: %w", err)
	}
	cache.Add(requestURL, cacheData)

	return data, nil
}
