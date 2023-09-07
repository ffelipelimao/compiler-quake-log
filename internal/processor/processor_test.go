package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateOutput(t *testing.T) {
	testCases := []struct {
		description string
		input       []string
		expected    *Output
	}{
		{
			description: "Empty input",
			input:       []string{},
			expected: &Output{
				Game: []Game{},
			},
		},
		{
			description: "Single game",
			input: []string{
				"InitGame:",
				"  1:08 Kill: 1 2 3: player1 killed player2 by MOD_WEAPON",
				"  1:08 Kill: 4 5 6: player3 killed player1 by MOD_WEAPON",
				"  1:08 Kill: 7 8 9: <world> killed player2 by MOD_WEAPON",
				"ShutdownGame:",
			},
			expected: &Output{
				Game: []Game{
					{
						TotalKills: 2,
						Players:    []string{"player1", "player2", "player3"},
						Kills: map[string]int{
							"player1": 1,
							"player2": -1,
							"player3": 1,
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			p := New()
			output := p.CreateOutput(tc.input)

			assert.Equal(t, len(tc.expected.Game), len(output.Game))

			for i, expectedGame := range tc.expected.Game {
				actualGame := output.Game[i]

				assert.Equal(t, expectedGame.TotalKills, actualGame.TotalKills)
				assert.ElementsMatch(t, expectedGame.Players, actualGame.Players)

				for player, expectedKills := range expectedGame.Kills {
					actualKills := actualGame.Kills[player]
					assert.Equal(t, expectedKills, actualKills)
				}
			}
		})
	}
}
