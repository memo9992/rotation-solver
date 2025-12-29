package main

import (
	"fmt"
	"time"

	"github.com/memo9992/rotation-solver/pkg/ff"
	"github.com/memo9992/rotation-solver/pkg/game"
	"github.com/memo9992/rotation-solver/pkg/solver"
)

func main() {
	r := solver.NewRotation(
		time.Millisecond*2500,
		time.Millisecond*1000,
	)
	rotation := []game.Ability{
		ff.HolySpirit,
		ff.FastBlade,
		ff.FightOrFlight,
		ff.RiotBlade,
		// ff.Requiescat,
	}

	for _, a := range rotation {
		r.AddEvent(a)
	}

	results := r.Calculate()

	for damageType, damageValue := range results.TotalDamage {
		fmt.Printf("Total %s Damage: %d\n", damageType, damageValue)
	}

	totalTime := results.TotalTime
	fmt.Printf("Total time: %s", totalTime.String())
}
