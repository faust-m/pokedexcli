package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"

	"github.com/faust-m/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

type config struct {
	next     *url.URL
	previous *url.URL
}

var cmds map[string]cliCommand
var pokedex map[string]pokeapi.Pokemon

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
		"explore": {
			name:        "explore",
			description: "List Pokemon found in a given area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "List attributes for a caught Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List Pokemon names in your Pokedex",
			callback:    commandPokedex,
		},
	}

	pokedex = map[string]pokeapi.Pokemon{}
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
				fmtInput := cleanInput(text)
				value, ok := cmds[fmtInput[0]]
				if !ok {
					fmt.Println("Unknown command")
					continue
				}
				if err := value.callback(&cfg, fmtInput[1:]...); err != nil {
					fmt.Println("Error:", err)
				}
			}
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(*config, ...string) error {
	if _, err := fmt.Println("Closing the Pokedex... Goodbye!"); err != nil {
		return fmt.Errorf("error in commandExit: %w", err)
	}
	os.Exit(0)
	return nil
}

func commandHelp(*config, ...string) error {
	heading :=
		`Welcome to the Pokedex!
Usage:

`
	if _, err := fmt.Print(heading); err != nil {
		return fmt.Errorf("error printing heading in commandHelp: %w", err)
	}
	for _, v := range cmds {
		if _, err := fmt.Printf("%-10s %s\n", v.name+":", v.description); err != nil {
			return fmt.Errorf("error printing command in commandHelp: %w", err)
		}
	}
	return nil
}

func commandMap(cfg *config, args ...string) error {
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

func commandMapb(cfg *config, args ...string) error {
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

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no area specified to explore")
	}
	fmt.Printf("Exploring %s...\n", args[0])
	exploreData, err := pokeapi.ExploreArea(fmt.Sprintf("%s%s/%s", pokeapi.BaseURL, pokeapi.LocationAreaEP, args[0]))
	if err != nil {
		return fmt.Errorf("error exploring %s", args[0])
	}
	if len(exploreData.PokemonEncounters) > 0 {
		fmt.Println("Found Pokemon:")
		for _, pokemon := range exploreData.PokemonEncounters {
			fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
		}
	}
	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no Pokemon specified to catch")
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])
	pokemonData, err := pokeapi.GetPokemonData(fmt.Sprintf("%s%s/%s", pokeapi.BaseURL, pokeapi.PokemonEP, args[0]))
	if err != nil {
		return fmt.Errorf("error getting Pokemon data: %w", err)
	}
	if rand.Intn(pokemonData.BaseExperience) <= 40 {
		pokedex[pokemonData.Name] = pokemonData
		fmt.Printf("%s was caught!\n", pokemonData.Name)
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped!\n", pokemonData.Name)
	}
	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no Pokemon specified to inspect")
	}
	data, ok := pokedex[args[0]]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\nHeight: %v\nWeight: %v\nStats:\n", data.Name, data.Height, data.Weight)
	for _, stat := range data.Stats {
		fmt.Printf(" -%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range data.Types {
		fmt.Printf(" - %s\n", t.Type.Name)
	}
	return nil
}

func commandPokedex(*config, ...string) error {
	if len(pokedex) == 0 {
		fmt.Println("you have no pokemon in your pokedex")
	} else {
		fmt.Println("Pokedex:")
		for k := range pokedex {
			fmt.Printf(" - %s\n", pokedex[k].Name)
		}
	}
	return nil
}
