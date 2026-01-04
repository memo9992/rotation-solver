package solver

import (
	"fmt"

	"github.com/memo9992/rotation-solver/pkg/game"
)

// Solve a rotation to determine the total damage output of a series of abilities, representing a rotation.
func Solve(rotation []game.Ability) {
	// TODO: Calculate time of abilities and durations of buffs.
	dmgInstances := []game.DamageInstance{}
	c := game.Conditions{}

	for _, r := range rotation {
		damage, conditionMods := r.Activate(c)
		dmgInstances = append(dmgInstances, damage...)
		for _, mods := range conditionMods {
			c.WithAction(mods)
		}
	}

	for _, d := range dmgInstances {
		fmt.Printf("Damage (%v): %v\n", d.Source, d.Calculate())
	}
}
