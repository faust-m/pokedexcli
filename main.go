package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cmds := initCommands()
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			if text := scanner.Text(); len(text) > 0 {
				value, ok := cmds[cleanInput(text)[0]]
				if !ok {
					fmt.Println("Unknown command")
					continue
				}
				value.callback()
			}
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func initCommands() map[string]cliCommand {
	cmds := map[string]cliCommand{
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
	}
	return cmds
}

func commandExit() error {
	if _, err := fmt.Println("Closing the Pokedex... Goodbye!"); err != nil {
		return fmt.Errorf("error in commandExit: %w", err)
	}
	os.Exit(0)
	return nil
}

func commandHelp() error {
	heading :=
		`Welcome to the Pokedex!
Usage:

`
	if _, err := fmt.Print(heading); err != nil {
		return fmt.Errorf("error printing heading in commandHelp: %w", err)
	}
	cmds := initCommands()
	for _, v := range cmds {
		if _, err := fmt.Printf("%s: %s\n", v.name, v.description); err != nil {
			return fmt.Errorf("error printing command in commandHelp: %w", err)
		}
	}
	return nil
}
