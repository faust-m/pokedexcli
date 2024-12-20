package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetLocationAreas(requestURL string) (LocationArea, error) {
	res, err := http.Get(requestURL)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error getting resource: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var data LocationArea
	err = decoder.Decode(&data)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error decoding response: %w", err)
	}

	return data, nil
}
