package main

import (
	"time"

	"github.com/memo9992/rotation-solver/pkg/game"
	"github.com/memo9992/rotation-solver/pkg/solver"
)

func main() {
	comboStarter := game.Ability{
		ID:         "comboStarter",
		Base:       100,
		DoesDamage: true,
		ConditionModifiers: []game.ConditionAction{
			{
				Action: game.ADD,
				Condition: game.Condition{
					ID:       "comboStarterCombo",
					Duration: time.Second * 30,
				},
			},
		},
	}
	comboFollower := game.Ability{
		ID:         "comboFollower",
		Base:       100,
		DoesDamage: true,
		DamageModifiers: []game.DamageModifier{
			{
				Action:              game.ADDITIVE,
				Value:               150,
				RequiredConditionID: "comboStarterCombo",
			},
		},
		ConditionModifiers: []game.ConditionAction{
			{
				Action: game.REMOVE,
				Condition: game.Condition{
					ID: "comboStarterCombo",
				},
			},
		},
	}

	buff := game.Ability{
		ID:         "buff",
		DoesDamage: false,
		ConditionModifiers: []game.ConditionAction{
			{
				Action: game.ADD,
				Condition: game.Condition{
					ID:       "buff",
					Duration: time.Second * 30,
					Modifiers: []game.DamageModifier{
						{
							Source: "buff",
							Action: game.INCREASE,
							Value:  10,
						},
					},
				},
			},
		},
	}

	rotation := []game.Ability{
		comboStarter,
		buff,
		comboFollower,
		comboFollower,
	}

	solver.Solve(rotation)
}
