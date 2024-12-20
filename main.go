package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
	//"github.com/faust-m/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	next     url.URL
	previous url.URL
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
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			if text := scanner.Text(); len(text) > 0 {
				value, ok := cmds[cleanInput(text)[0]]
				if !ok {
					fmt.Println("Unknown command")
					continue
				}
				switch value.name {
				case "map":
					value.callback(&config{})
				default:
					value.callback(&config{})
				}
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
	return nil
}
