package solver

import (
	"fmt"
	"time"

	"github.com/memo9992/rotation-solver/pkg/game"
)

type Rotation struct {
	gcdDuration  time.Duration
	ogcdDuration time.Duration
	events       []game.Ability
	ongoingBuffs *game.Buffs
}

func NewRotation(gcdDuration time.Duration, ogcdDuration time.Duration) *Rotation {
	return &Rotation{
		gcdDuration:  gcdDuration,
		ogcdDuration: ogcdDuration,
		ongoingBuffs: &game.Buffs{},
	}
}

func (r *Rotation) AddEvent(a game.Ability) {
	r.events = append(r.events, a)
}

type Result struct {
	TotalDamage map[game.DamageType]int64
	TotalTime   time.Duration
}

func (r *Rotation) Calculate() Result {
	totalDamage := make(map[game.DamageType]int64)
	var now time.Time
	var gcdTimeRemaining time.Duration

	for _, ability := range r.events {
		fmt.Println(now)
		for _, b := range r.ongoingBuffs.State() {
			fmt.Printf("\t%s: %s (%v)\n", b.Name, b.TimeStart, b.Duration)
		}
		act := ability.OnActivate(now, r.ongoingBuffs)
		r.ongoingBuffs.Update(now)
		r.ongoingBuffs.Apply(act)
		ability = ability.AfterBuffs(ability, r.ongoingBuffs)

		// TODO: Does not work ATM
		var damageDoneMultiplier int64 = 100

		for _, buff := range r.ongoingBuffs.State() {
			if mult, ok := buff.Effects[game.BuffEffectDamageMultiplier]; ok {
				damageDoneMultiplier += mult
			}
		}

		totalDamage[ability.DamageType] += ability.Damage * damageDoneMultiplier / 100

		if ability.IsGCD() {
			now = now.Add(max(r.gcdDuration, ability.CastTime))
			gcdTimeRemaining = r.gcdDuration
		} else if ability.IsOGCD() {
			actualCastTime := max(ability.CastTime, r.ogcdDuration)

			if actualCastTime <= gcdTimeRemaining {
				gcdTimeRemaining -= actualCastTime
			} else {
				now = now.Add(actualCastTime - gcdTimeRemaining)
				gcdTimeRemaining = 0
			}
		}
	}

	totalTime := now.Sub(time.Time{})

	return Result{
		TotalDamage: totalDamage,
		TotalTime:   totalTime,
	}
}
