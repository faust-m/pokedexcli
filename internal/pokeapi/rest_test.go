package pokeapi

import (
	"testing"
)

func TestInvalidResource(t *testing.T) {
	url := "https://pokeapi.co/api/v2/abcd/invalid"
	_, err := GetLocationAreas(url)
	if err == nil {
		t.Errorf("invalid request did not err")
		t.Fail()
	}
}

func TestValidNonCacheResource(t *testing.T) {
	url := BaseURL + LocationAreaEP
	locationData, _ := GetLocationAreas(url)
	if len(locationData.Results) == 0 {
		t.Errorf("valid request has zero result length")
		t.Fail()
	}
}

func TestValidCacheResource(t *testing.T) {
	url := BaseURL + LocationAreaEP
	GetLocationAreas(url)
	if _, found := cache.Get(url); !found {
		t.Errorf("request was not cached")
		t.Fail()
	}
}
