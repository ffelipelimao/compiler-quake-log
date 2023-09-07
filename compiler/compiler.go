package compiler

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

type Compiler struct {
	Path        string
	Rows        []string
	CurrentGame *Game
}

type Event struct {
	Game []Game `json:"game"`
}

type Game struct {
	TotalKills int            `json:"total_kills"`
	Players    []string       `json:"players"`
	Kills      map[string]int `json:"kills"`
}

func New(path string) *Compiler {
	return &Compiler{
		Path: path,
	}
}

func (c *Compiler) LoadRows() error {
	var rows []string

	file, err := os.Open(c.Path)
	if err != nil {
		return fmt.Errorf("error to open file %s", err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		rows = append(rows, row)
	}

	c.Rows = rows

	return nil
}

func (c *Compiler) Process() *Event {
	isGameLine := false
	event := &Event{}

	for _, row := range c.Rows {
		row = strings.TrimSpace(row)

		if strings.Contains(row, "InitGame:") {
			isGameLine = true
			c.CurrentGame = &Game{
				TotalKills: 0,
				Players:    []string{},
				Kills:      make(map[string]int),
			}
			continue
		}

		if strings.Contains(row, "ShutdownGame:") {
			isGameLine = false
			if c.CurrentGame != nil {
				event.Game = append(event.Game, *c.CurrentGame)
			}
			continue
		}

		if isGameLine {
			c.processGame(row)
		}
	}

	return event
}

func (c *Compiler) processGame(row string) {
	row = strings.TrimSpace(row)
	if strings.Contains(row, "Kill") {

		killedRegex := regexp.MustCompile(`[^:]+:[^:]+:[^:]+:\s+[^:]+\s+killed\s+(.*?)\s+by`)

		killerRegex := regexp.MustCompile(`[^:]+:[^:]+:[^:]+:\s*(.*?)\s+killed`)

		killedName := killedRegex.FindStringSubmatch(row)

		killerName := killerRegex.FindStringSubmatch(row)

		if len(killerName) >= 2 && len(killedName) >= 2 {
			killer := killerName[1]
			killed := killedName[1]

			if killer != "<world>" {

				c.CurrentGame.TotalKills++
				c.CurrentGame.Kills[killer]++

				if !slices.Contains(c.CurrentGame.Players, killer) {
					c.CurrentGame.Players = append(c.CurrentGame.Players, killer)
				}
			} else {
				c.CurrentGame.Kills[killed]--
				if !slices.Contains(c.CurrentGame.Players, killed) {
					c.CurrentGame.Players = append(c.CurrentGame.Players, killed)
				}
			}
		}
	}
}
