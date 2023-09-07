package processor

import (
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

type Processor struct {
	CurrentGame *Game
}

type Output struct {
	Game []Game `json:"game"`
}

type Game struct {
	TotalKills int            `json:"total_kills"`
	Players    []string       `json:"players"`
	Kills      map[string]int `json:"kills"`
}

func New() *Processor {
	return &Processor{}
}

func (p *Processor) CreateOutput(rows []string) *Output {
	isGameLine := false
	output := &Output{}

	for _, row := range rows {
		row = strings.TrimSpace(row)

		if strings.Contains(row, "InitGame:") {
			isGameLine = true
			p.CurrentGame = &Game{
				TotalKills: 0,
				Players:    []string{},
				Kills:      make(map[string]int),
			}
			continue
		}

		if strings.Contains(row, "ShutdownGame:") {
			isGameLine = false
			if p.CurrentGame != nil {
				output.Game = append(output.Game, *p.CurrentGame)
			}
			continue
		}

		if isGameLine {
			p.processGame(row)
		}
	}

	return output
}

func (c *Processor) processGame(row string) {
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
