package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/faust-m/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	next     *url.URL
	previous *url.URL
}

var cmds map[string]cliCommand

func init() {
	cmds = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 location areas",
			callback:    commandMapb,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := config{
		next:     &url.URL{},
		previous: &url.URL{},
	}
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			if text := scanner.Text(); len(text) > 0 {
				value, ok := cmds[cleanInput(text)[0]]
				if !ok {
					fmt.Println("Unknown command")
					continue
				}
				value.callback(&cfg)
			}
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(*config) error {
	if _, err := fmt.Println("Closing the Pokedex... Goodbye!"); err != nil {
		return fmt.Errorf("error in commandExit: %w", err)
	}
	os.Exit(0)
	return nil
}

func commandHelp(*config) error {
	heading :=
		`Welcome to the Pokedex!
Usage:

`
	if _, err := fmt.Print(heading); err != nil {
		return fmt.Errorf("error printing heading in commandHelp: %w", err)
	}
	for _, v := range cmds {
		if _, err := fmt.Printf("%s:\t%s\n", v.name, v.description); err != nil {
			return fmt.Errorf("error printing command in commandHelp: %w", err)
		}
	}
	return nil
}

func commandMap(cfg *config) error {
	var err error
	if cfg.next == nil {
		fmt.Println("You're on the last page!")
		return nil
	} else if cfg.next.String() == "" {
		cfg.next, err = url.Parse(pokeapi.BaseURL + pokeapi.LocationAreaEP)
		if err != nil {
			return fmt.Errorf("error parsing URL in commandMap: %w", err)
		}
		q := url.Values{}
		q.Add(pokeapi.OffsetKey, "0")
		q.Add(pokeapi.LimitKey, "20")
		cfg.next.RawQuery = q.Encode()
	}

	locationAreas, err := pokeapi.GetLocationAreas(cfg.next.String())
	if err != nil {
		return fmt.Errorf("error getting next location areas: %w", err)
	}
	err = updateConfig(cfg, locationAreas)
	if err != nil {
		return fmt.Errorf("error updating config: %w", err)
	}

	for _, result := range locationAreas.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandMapb(cfg *config) error {
	var err error
	if cfg.previous == nil {
		fmt.Println("You're on the first page!")
		return nil
	} else if cfg.previous.String() == "" {
		cfg.previous, err = url.Parse(pokeapi.BaseURL + pokeapi.LocationAreaEP)
		if err != nil {
			return fmt.Errorf("error parsing URL in commandMapb: %w", err)
		}
		q := url.Values{}
		q.Add(pokeapi.OffsetKey, "0")
		q.Add(pokeapi.LimitKey, "20")
		cfg.next.RawQuery = q.Encode()
	}

	locationAreas, err := pokeapi.GetLocationAreas(cfg.previous.String())
	if err != nil {
		return fmt.Errorf("error getting next location areas: %w", err)
	}
	err = updateConfig(cfg, locationAreas)
	if err != nil {
		return fmt.Errorf("error updating config: %w", err)
	}

	for _, result := range locationAreas.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func updateConfig(cfg *config, locationAreas pokeapi.LocationArea) error {
	var err error
	if locationAreas.Previous != nil {
		cfg.previous, err = url.Parse(*locationAreas.Previous)
		if err != nil {
			return fmt.Errorf("error parsing previous URL: %w", err)
		}
	} else {
		cfg.previous = nil
	}
	if locationAreas.Next != nil {
		cfg.next, err = url.Parse(*locationAreas.Next)
		if err != nil {
			return fmt.Errorf("error parsing next URL: %w", err)
		}
	} else {
		cfg.next = nil
	}
	return nil
}
